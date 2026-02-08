import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

export interface HealthStatus {
  healthy: boolean;
}

export interface AboutInfo {
  status: string;
  message: string;
}

// Demo data interfaces
export interface DemoData {
  id: number;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface DemoDataRequest {
  content: string;
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  // TMP: Use environment variables via environment.ts and environment.prod.ts
  // Angular CLI supports file replacements for different configurations
  private baseUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  checkHealth(): Observable<HealthStatus> {
    return this.http.get<{status: string}>(`${this.baseUrl}/health`).pipe(
      map(response => ({ healthy: response.status === 'ok' })),
      catchError(error => {
        console.error('Health check failed:', error);
        return of({ healthy: false });
      })
    );
  }

  getInfo(): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/info`).pipe(
      catchError(error => {
        console.error('Info fetch failed:', error);
        return of(null);
      })
    );
  }

  getAbout(): Observable<AboutInfo> {
    return of({
      status: 'healthy',
      message: 'K8s Fullstack Blueprint'
    });
  }

  // Demo data API methods
  getAllDemoData(): Observable<DemoData[]> {
    return this.http.get<DemoData[]>(`${this.baseUrl}/api/demo-data`).pipe(
      catchError(error => {
        console.error('Failed to fetch demo data:', error);
        return of([]);
      })
    );
  }

  upsertDemoData(content: string): Observable<DemoData> {
    return this.http.post<DemoData>(`${this.baseUrl}/api/demo-data`, { content }).pipe(
      catchError(error => {
        console.error('Failed to upsert demo data:', error);
        throw error;
      })
    );
  }
}
