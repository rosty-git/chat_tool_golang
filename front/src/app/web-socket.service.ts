import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Injectable({
  providedIn: 'root',
})
export class WebSocketService {
  private socket$!: WebSocketSubject<unknown>;

  connect(url: string): void {
    this.socket$ = webSocket(url);
  }

  sendMessage(message: unknown): void {
    if (this.socket$) {
      this.socket$.next(message);
    } else {
      console.error('WebSocket connection is not established.');
    }
  }

  getMessages(): Observable<unknown> {
    return this.socket$.asObservable();
  }

  close(): void {
    if (this.socket$) {
      this.socket$.complete();
    }
  }
}
