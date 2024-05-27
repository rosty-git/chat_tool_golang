import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/internal/Observable';

import { GlobalVariable } from '../global';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  constructor(private http: HttpClient) {}

  post(path: string, data: unknown): Observable<unknown> {
    return this.http.post(`${GlobalVariable.BASE_API_URL}${path}`, data, {
      withCredentials: true,
    });
  }

  put(path: string, data: unknown): Observable<unknown> {
    return this.http.put(`${GlobalVariable.BASE_API_URL}${path}`, data, {
      withCredentials: true,
    });
  }

  get(path: string, params?: HttpParams): Observable<unknown> {
    return this.http.get(`${GlobalVariable.BASE_API_URL}${path}`, {
      withCredentials: true,
      params,
    });
  }
}
