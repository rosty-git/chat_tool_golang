import { Component } from '@angular/core';

import { DirectChannelsComponent } from '../directChannels/directChannels.component';
import { OpenChannelsComponent } from '../openChannels/openChannels.component';

// export type Channel = {
//   id: string;
//   name: string;
// };

// export type ChannelsResp = {
//   channels: Channel[];
// };

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [OpenChannelsComponent, DirectChannelsComponent],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss',
})
export class SidebarComponent {
  // constructor(private api: ApiService) {}

  // directChannels: Channel[] = [];

  // openChannels: Channel[] = [];

  // ngOnInit(): void {
  //   const directParams = new HttpParams().append('channelType', 'D');

  //   this.api.get('/v1/api/channels', directParams).subscribe({
  //     next: (response) => {
  //       console.log('Get Channels', response);

  //       this.directChannels = (response as ChannelsResp).channels;
  //     },

  //     error: (err: unknown) => {
  //       console.error('error', err);
  //     },
  //   });

  //   const openParams = new HttpParams().append('channelType', 'O');

  //   this.api.get('/v1/api/channels', openParams).subscribe({
  //     next: (response) => {
  //       console.log('Get Contacts', response);

  //       this.openChannels = (response as ChannelsResp).channels;
  //     },

  //     error: (err) => {
  //       console.error('error', err);
  //     },
  //   });
  // }
}
