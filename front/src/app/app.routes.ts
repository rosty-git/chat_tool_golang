import { Routes } from '@angular/router';

import { LoginComponent } from './login/login.component';
import { MessengerComponent } from './messenger/messenger.component';

export const routes: Routes = [
  { path: '', redirectTo: 'messenger', pathMatch: 'full' },
  { path: 'messenger', component: MessengerComponent },
  { path: 'login', component: LoginComponent },
  { path: '**', component: MessengerComponent },
];
