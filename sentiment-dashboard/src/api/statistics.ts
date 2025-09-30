import axios from 'axios';

export type SentimentPoint = {
  date: string; // ISO
  positive: number;
  negative: number;
  neutral: number;
};

export type ProductStats = {
  name: string;
  timeline: SentimentPoint[];
};

export type StatisticsResponse = {
  products: ProductStats[];
};

const api = axios.create({
  baseURL: (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080',
});

export async function fetchStatistics(): Promise<StatisticsResponse> {
  const resp = await api.get<StatisticsResponse>('/statistics');
  return resp.data;
}


