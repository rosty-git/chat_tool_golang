import { Component, HostListener } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

import { ApiService } from '../api.service';
import { DataService } from '../data.service';
import { DirectChannelComponent } from '../direct-channel/direct-channel.component';
import { OpenChannelComponent } from '../open-channel/open-channel.component';

type Channel = {
  channelType: 'D' | 'O';
  id: string;
  name: string;
};

type SearchResults = {
  results: Channel[];
};

@Component({
  selector: 'app-search-channels',
  standalone: true,
  imports: [ReactiveFormsModule, OpenChannelComponent, DirectChannelComponent],
  templateUrl: './search-channels.component.html',
  styleUrl: './search-channels.component.scss',
})
export class SearchChannelsComponent {
  isOpen = false;

  channels: Channel[] = [];

  searchForm = new FormGroup({
    text: new FormControl(''),
  });

  constructor(
    private dataService: DataService,
    private api: ApiService,
  ) {
    this.dataService.showChannelSearchModal$.subscribe((value) => {
      this.isOpen = value;
    });

    this.searchForm.valueChanges.subscribe((value) => {
      this.api.get(`/v1/api/channels/search/${value.text}`).then((resp) => {
        console.log(resp);

        this.channels = (resp as SearchResults).results;
      });
    });
  }

  closeModal() {
    this.searchForm.reset();
    this.dataService.closeChannelSearchModal();
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (this.isOpen && target.classList.contains('modal-backdrop')) {
      this.closeModal();
    }
  }
}
