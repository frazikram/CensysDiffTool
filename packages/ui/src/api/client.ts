export interface Snapshot {
  id?: string;
  ip: string;
  timestamp: string;
  services: Service[];
}

export interface Service {
  port: number;
  protocol: string;
  status?: number;
  software?: {
    vendor: string;
    product: string;
    version?: string;
  };
  vulnerabilities?: string[];
  tls?: {
    version?: string;
    cipher?: string;
  };
}

export interface DiffResult {
  added: Service[];
  removed: Service[];
  changed: { before: Service; after: Service }[];
}

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

class ApiClient {
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      headers: { "Content-Type": "application/json", ...options?.headers },
      ...options,
    });
    if (!response.ok) {
      throw new Error(await response.text());
    }
    return response.json();
  }

  uploadSnapshot = async (file: File) =>
    this.request<Snapshot>("/snapshots", {
      method: "POST",
      body: await file.text(),
    });

  getHistory = async (ip: string) =>
    this.request<Snapshot[]>(`/hosts/${ip}/snapshots`);

  getDiff = async (snap1: string, snap2: string) =>
    this.request<DiffResult>(`/diff?snap1=${snap1}&snap2=${snap2}`);

  getHosts = async () => this.request<string[]>("/all-hosts");
}

export const apiClient = new ApiClient();
