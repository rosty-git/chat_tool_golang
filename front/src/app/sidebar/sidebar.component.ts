import { Component, OnInit } from '@angular/core';

import { ApiService } from '../api.service';
import { ChannelsComponent } from '../channels/channels.component';
import { ContactsComponent } from '../contacts/contacts.component';

export type Channel = {
  id: number;
  name: string;
};

export type Contact = {
  id: number;
  name: string;
};

export type ChannelsResp = {
  channels: Channel[];
};

export type ContactsResp = {
  contacts: Contact[];
};

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [ChannelsComponent, ContactsComponent],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss',
})
export class SidebarComponent implements OnInit {
  constructor(private api: ApiService) {}

  channels: Channel[] = [];

  contacts: Contact[] = [];

  ngOnInit(): void {
    this.api.get('/v1/api/channels').subscribe({
      next: (response) => {
        console.log('Get Channels', response);

        this.channels = (response as ChannelsResp).channels;
      },

      error: (err: unknown) => {
        console.error('error', err);
      },
    });

    this.api.get('/v1/api/contacts').subscribe({
      next: (response) => {
        console.log('Get Contacts', response);

        this.contacts = (response as ContactsResp).contacts;
      },

      error: (err) => {
        console.error('error', err);
      },
    });
  }
}
