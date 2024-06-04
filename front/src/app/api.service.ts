import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { GlobalVariable } from '../global';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  constructor(private http: HttpClient) {}

  post(path: string, data: unknown): Promise<unknown> {
    return new Promise((resole, reject) => {
      this.http
        .post(`${GlobalVariable.BASE_API_URL}${path}`, data, {
          withCredentials: true,
        })
        .subscribe({
          next: (response) => {
            resole(response);
          },
          error: (err: unknown) => {
            console.error('error', err);

            reject(err);
          },
        });
    });
  }

  put<T>(path: string, data: unknown): Promise<T> {
    return new Promise((resole, reject) => {
      this.http
        .put(`${GlobalVariable.BASE_API_URL}${path}`, data, {
          withCredentials: true,
        })
        .subscribe({
          next: (response) => {
            resole(response as T);
          },
          error: (err: unknown) => {
            console.error('error', err);

            reject(err);
          },
        });
    });
  }

  get(path: string, params?: HttpParams): Promise<unknown> {
    return new Promise((resole, reject) => {
      this.http
        .get(`${GlobalVariable.BASE_API_URL}${path}`, {
          withCredentials: true,
          params,
        })
        .subscribe({
          next: (response) => {
            resole(response);
          },
          error: (err: unknown) => {
            console.error('error', err);

            reject(err);
          },
        });
    });
  }
}
