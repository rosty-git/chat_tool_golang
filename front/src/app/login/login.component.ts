// import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { retry } from 'rxjs';
// import { WebSocketService } from '../web-socket.service';
import { webSocket } from 'rxjs/webSocket';

import { GlobalVariable } from '../../global';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent {
  constructor(
    private api: ApiService,
    // private webSocket: WebSocketService,
  ) {}

  loginForm = new FormGroup({
    login: new FormControl(''),
    password: new FormControl(''),
  });

  submit() {
    console.log('SUBMIT');

    this.api
      .post('/v1/auth/login', {
        email: this.loginForm.value.login,
        password: this.loginForm.value.password,
      })
      .subscribe({
        next: (response: unknown) => {
          console.log('auth ok', response);
          this.loginForm.reset();
        },

        error: (err) => {
          console.error('auth error', err);
        },

        complete() {
          console.log('is completed');
        },
      });
  }

  showMessages(): void {
    const subject = webSocket(`${GlobalVariable.BASE_WS_URL}/v1/ws/`);

    subject.subscribe({
      next: (msg) => console.log('message received', msg), // Called whenever there is a message from the server.
      error: (err) => {
        console.log(err);
      }, // Called if at any point WebSocket API signals some kind of error.
      complete: () => console.log('complete'), // Called when connection is closed (for whatever reason).
    });
  }
}
