export interface User {
    full_name: string;
    id: string;
    email: string;
    imageUrl?: string;
    role: 'owner' | 'manager' | 'staff';
    business_name: string;
    telephone: string;
    location: string;
  }
