import axios from 'axios';

const API_URL = "http://localhost:3000";

export interface Subject {
    id: number;
    title: string;
    slug: string;
    color: string;
    language: string;
    iconName: string;
}

export interface SubjectSession {
    id: number;
    subjectId: number;
    title: string;
    orderIndex: number;
}

export interface Lesson {
    id: number;
    sessionId: number;
    title: string;
    contentMd: string;
    orderIndex: number;
}

export const courseService = {
    getSubjects: async (): Promise<Subject[]> => {
        const response = await axios.get(`${API_URL}/subjects`);
        return response.data;
    },
    getSubjectSessions: async (subjectId?: number): Promise<SubjectSession[]> => {
        const url = subjectId ? `${API_URL}/sessions?subjectId=${subjectId}` : `${API_URL}/sessions`;
        const response = await axios.get(url);
        return response.data;
    },
    getLessonsBySession: async (sessionId?: number): Promise<Lesson[]> => {
        const url = sessionId ? `${API_URL}/lessons?sessionId=${sessionId}` : `${API_URL}/lessons`;
        const response = await axios.get(url);
        return response.data;
    }
};
