// import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';

// import { WebSocketService } from '../web-socket.service';
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
    private router: Router,
  ) {}

  loginForm = new FormGroup({
    login: new FormControl(''),
    password: new FormControl(''),
  });

  submit() {
    this.api
      .post('/v1/auth/login', {
        email: this.loginForm.value.login,
        password: this.loginForm.value.password,
      })
      .subscribe({
        next: (response: unknown) => {
          console.log('auth ok', response);
          this.loginForm.reset();
          this.router.navigate(['messenger']);
        },

        error: (err) => {
          console.error('auth error', err);
        },
      });
  }
}
