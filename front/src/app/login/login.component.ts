import { Component } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';

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
      .then(() => {
        this.loginForm.reset();
        this.router.navigate(['messenger']);
      });
  }
}
