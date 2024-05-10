import { Component } from '@angular/core';
import {FormControl, FormGroup, ReactiveFormsModule} from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { GlobalVariable } from '../../global';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {
  constructor(private http: HttpClient) {}

  loginForm = new FormGroup({
    login: new FormControl(''),
    password: new FormControl(''),
  });

  submit() {
    console.log("SUBMIT")

    this.http.post(GlobalVariable.BASE_API_URL + '/v1/auth/login', 
    { email: this.loginForm.value.login, password: this.loginForm.value.password }).subscribe(response => {
      console.log('Login successful', response);
      this.http.get(GlobalVariable.BASE_API_URL + '/v1/api/messages', {withCredentials: true}).subscribe(response => {
        console.log('get resp', response);
      }, error => {
        console.error('get error', error);
      });
    }, error => {
      console.error('Login failed', error);
    });

    // this.loginForm.reset();
  }

  showMessages() {
    this.http.get(GlobalVariable.BASE_API_URL + '/v1/api/messages', {withCredentials: true}).subscribe(response => {
      console.log('get resp', response);
    }, error => {
      console.error('get error', error);
    });
  }
}
