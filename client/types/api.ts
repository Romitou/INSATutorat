export interface Campaign {
    id: number;
    semester: number;
    schoolYear: string;
    startDate: string;
    endDate: string;
    registrationStatus: string;
    registrationStartDate: string;
    registrationEndDate: string;
}

export interface SemesterAvailability {
    id: number;
    campaign: Campaign;
    campaignId: number;
    user: User;
    userId: number;
    availabilityJson: string;
}

export interface Subject {
    id: number;
    semester: number;
    shortName: string;
    name: string;
}

export interface TuteeRegistration {
    id: number;
    tutee: User;
    tuteeId: number;
    subject: Subject;
    subjectId: number;
    tutorSubjectId?: number | null;
    totalHours: number;
}

export interface TutoringHour {
    id?: number | null;
    tuteeId?: number | null;
    startDate: string;
    endDate: string;
}

export interface TutoringLesson {
    id: number;
    content: string;
    startDate: string;
    endDate: string;
}

export interface TutorSubject {
    id: number;
    tutor: User;
    tutorId: number;
    subject: Subject;
    subjectId: number;
    maxTutees: number;
    totalHours: number;
}

export interface User {
    id: number;
    firstName: string;
    lastName: string;
    mail: string;
    schoolYear: string;
    groups: string[];
    isTutor: boolean;
    isTutee: boolean;
    isAdmin: boolean;
}

export interface TuteeAssignment {
    id: number
    tutee: User
    tuteeId: number
    subjectId: number
    tutorSubjectId: number | null
    totalHours: number
}

export type Slots = {
    [key: string]: number[]; // format YYYY-MM-DD
}

export const statuses = ['AVAILABLE', 'OCCUPIED'] as const;
export type SlotStatus = typeof statuses[number];