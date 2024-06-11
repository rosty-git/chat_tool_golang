import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-message-item',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './message-item.component.html',
  styleUrl: './message-item.component.scss',
})
export class MessageItemComponent {
  @Input() message: string = '';

  @Input() userName: string = '';

  @Input() id: string = '';

  @Input() created_at: string = '';

  @Input() offline: boolean = false;
}
