import { Injectable } from '@angular/core';

import { GlobalVariable } from '../global';

@Injectable({
  providedIn: 'root',
})
export class WebSocketService {
  private webSocket: Socket;

  constructor() {
    this.webSocket = new Socket({
      url: 'localhost:8080/',
      options: {
        path: '/v1/ws',
        withCredentials: true,
        reconnectionDelay: 5000,
        autoConnect: false,
        // transports: ['websocket'],
      },
    });

    this.webSocket.connect();
  }

  connectSocket(message: unknown) {
    this.webSocket.emit('connect', message);
  }

  // this method is used to get response from server
  receiveStatus() {
    return this.webSocket.fromEvent('/get-response');
  }

  // this method is used to end web socket connection
  disconnectSocket() {
    this.webSocket.disconnect();
  }
}
