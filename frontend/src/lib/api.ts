// Mock API functions for SQL Sharding Management Tool
// In production, these would call the Go backend at localhost:8080

export interface Project {
  id: string;
  name: string;
  // connectionString: string;
  createdAt: string;
  shardsCount: number;
  status: 'active' | 'inactive';
}

export interface Shard {
  id: string;
  projectId: string;
  hostname: string;
  port: number;
  database: string;
  status: 'online' | 'offline' | 'maintenance';
  load: number;
  connections: number;
  lastPing: string;
}

export interface QueryResult {
  columns: string[];
  rows: any[][];
  executionTime: number;
  affectedRows?: number;
  error?: string;
}

// Mock data
const mockProjects: Project[] = [
  {
    id: '1',
    name: 'E-commerce DB',
    // connectionString: 'postgresql://user:pass@localhost:5432/ecommerce',
    createdAt: '2024-01-15T10:00:00Z',
    shardsCount: 4,
    status: 'active',
  },
  {
    id: '2',
    name: 'Analytics Platform',
    // connectionString: 'postgresql://user:pass@localhost:5432/analytics',
    createdAt: '2024-01-20T14:30:00Z',
    shardsCount: 6,
    status: 'active',
  },
  {
    id: '3',
    name: 'User Management',
    // connectionString: 'postgresql://user:pass@localhost:5432/users',
    createdAt: '2024-02-01T09:15:00Z',
    shardsCount: 2,
    status: 'inactive',
  },
];

const mockShards: Shard[] = [
  {
    id: 'shard-1',
    projectId: '1',
    hostname: 'shard1.example.com',
    port: 5432,
    database: 'ecommerce_shard_1',
    status: 'online',
    load: 45,
    connections: 12,
    lastPing: '2024-01-15T12:00:00Z',
  },
  {
    id: 'shard-2',
    projectId: '1',
    hostname: 'shard2.example.com',
    port: 5432,
    database: 'ecommerce_shard_2',
    status: 'online',
    load: 67,
    connections: 18,
    lastPing: '2024-01-15T12:00:00Z',
  },
  {
    id: 'shard-3',
    projectId: '1',
    hostname: 'shard3.example.com',
    port: 5432,
    database: 'ecommerce_shard_3',
    status: 'maintenance',
    load: 0,
    connections: 0,
    lastPing: '2024-01-15T11:45:00Z',
  },
];

// Simulate API delay
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// API Functions
export const getProjects = async (): Promise<Project[]> => {
  await delay(300);
  return [...mockProjects];
};

export const getProject = async (id: string): Promise<Project | null> => {
  await delay(200);
  return mockProjects.find(p => p.id === id) || null;
};

export const addProject = async (project: Omit<Project, 'id' | 'createdAt' | 'shardsCount' | 'status'>): Promise<Project> => {
  await delay(500);
  const newProject: Project = {
    ...project,
    id: Math.random().toString(36).substr(2, 9),
    createdAt: new Date().toISOString(),
    shardsCount: 0,
    status: 'active',
  };
  mockProjects.push(newProject);
  return newProject;
};


export const deleteProject = async (projectId: string): Promise<void> => {
  await delay(300);
  const index = mockProjects.findIndex(p => p.id === projectId);
  if (index !== -1) {
    mockProjects.splice(index, 1); // Remove from array in-place
  } else {
    throw new Error("Project not found");
  }
};


export const getShards = async (projectId: string): Promise<Shard[]> => {
  await delay(300);
  return mockShards.filter(shard => shard.projectId === projectId);
};

export const addShard = async (shard: Omit<Shard, 'id' | 'status' | 'load' | 'connections' | 'lastPing'>): Promise<Shard> => {
  await delay(400);
  const newShard: Shard = {
    ...shard,
    id: `shard-${Math.random().toString(36).substr(2, 9)}`,
    status: 'online',
    load: Math.floor(Math.random() * 100),
    connections: Math.floor(Math.random() * 50),
    lastPing: new Date().toISOString(),
  };
  mockShards.push(newShard);
  return newShard;
};

export const runQuery = async (projectId: string, query: string, shardId?: string): Promise<QueryResult> => {
  await delay(800);
  
  // Simulate different query results
  if (query.toLowerCase().includes('error')) {
    return {
      columns: [],
      rows: [],
      executionTime: 0,
      error: 'Syntax error in SQL query',
    };
  }

  if (query.toLowerCase().includes('select')) {
    return {
      columns: ['id', 'name', 'email', 'created_at'],
      rows: [
        [1, 'John Doe', 'john@example.com', '2024-01-15T10:00:00Z'],
        [2, 'Jane Smith', 'jane@example.com', '2024-01-16T11:30:00Z'],
        [3, 'Bob Johnson', 'bob@example.com', '2024-01-17T09:15:00Z'],
      ],
      executionTime: 0.045,
    };
  }

  if (query.toLowerCase().includes('insert') || query.toLowerCase().includes('update') || query.toLowerCase().includes('delete')) {
    return {
      columns: [],
      rows: [],
      executionTime: 0.032,
      affectedRows: Math.floor(Math.random() * 10) + 1,
    };
  }

  return {
    columns: ['result'],
    rows: [['Query executed successfully']],
    executionTime: 0.025,
  };
};

export const updateProjectSettings = async (projectId: string, settings: any): Promise<boolean> => {
  await delay(400);
  // Mock update
  return true;
};