import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

import { FrontFile } from '../data.service';
import { FileComponent } from '../file/file.component';

@Component({
  selector: 'app-message-item',
  standalone: true,
  templateUrl: './message-item.component.html',
  styleUrl: './message-item.component.scss',
  imports: [CommonModule, FileComponent]
})
export class MessageItemComponent {
  @Input() message: string = '';

  @Input() userName: string = '';

  @Input() id: string = '';

  @Input() created_at: string = '';

  @Input() offline: boolean = false;

  @Input() files: FrontFile[] = [];
}
