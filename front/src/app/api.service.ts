import { HttpClient } from '@angular/common/http';
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

  get(path: string): Observable<unknown> {
    return this.http.get(`${GlobalVariable.BASE_API_URL}${path}`, {
      withCredentials: true,
    });
  }
}
