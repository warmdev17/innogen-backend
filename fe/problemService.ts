import axios from 'axios';

const API_URL = "http://localhost:3000";

export interface Problem {
    id: number;
    lessonId: number;
    slug: string;
    acceptanceRate: number;
    title: string;
    difficulty: 'Easy' | 'Medium' | 'Hard';
    problemMd: string;
    template: string;
    timeLimitMs: number;
    memoryLimitMb: number;
    createdBy: number;
}

export interface TestCase {
    id: number;
    problemId: number;
    inputData: string;
    expectedOutput: string;
    isHidden: boolean;
    role: 'sample' | 'hidden' | 'edge_case';
}

export interface Submission {
    id: string; // UUID
    userId: number;
    problemId: number;
    language: string;
    code: string;
    status: 'Accepted' | 'Wrong Answer' | 'Time Limit Exceeded' | 'Compile Error' | 'Runtime Error';
    runtimeMs: number;
    memoryKb: number;
    errorMessage: string | null;
    passCount: number;
    totalTestcases: number;
}

export interface PistonExecutionResult {
    language: string;
    version: string;
    run: {
        stdout: string;
        stderr: string;
        code: number;
        signal: string | null;
        output: string;
        wall_time: number;
        memory: number;
    }
}

export const problemService = {
    getProblemsByLesson: async (lessonId?: number): Promise<Problem[]> => {
        const url = lessonId ? `${API_URL}/problems?lessonId=${lessonId}` : `${API_URL}/problems`;
        const response = await axios.get(url);
        return response.data;
    },
    getProblemById: async (id: number): Promise<Problem> => {
        const response = await axios.get(`${API_URL}/problems/${id}`);
        return response.data;
    },
    getProblemTestCases: async (problemId: number): Promise<TestCase[]> => {
        const response = await axios.get(`${API_URL}/test_cases?problemId=${problemId}`);
        return response.data;
    },
    submitCode: async (submission: Omit<Submission, 'id'>): Promise<Submission> => {
        const response = await axios.post(`${API_URL}/submissions`, {
            id: crypto.randomUUID(),
            ...submission
        });
        return response.data;
    },
    executeCode: async (language: string, _solution: string, _input: string): Promise<PistonExecutionResult> => {
        void _solution;
        void _input;
        // Mocking execution since backend is not available
        return new Promise((resolve) => {
            setTimeout(() => {
                // Determine mock validation logic depending on input
                // For simplicity, we assume successful compile and run
                // and return the expected mock data. 
                // We fake the expected output logic later in the slice since we don't have it here. This just mocks the engine.
                resolve({
                    language: language,
                    version: "mock",
                    run: {
                        stdout: "mock_output\n", // The caller slice will handle the actual mock matching logic to fake a success.
                        stderr: "",
                        code: 0,
                        signal: null,
                        output: "mock_output\n",
                        wall_time: Math.floor(Math.random() * 50) + 10,
                        memory: Math.floor(Math.random() * 3000) + 1024,
                    }
                });
            }, 800);
        });
    }
};
