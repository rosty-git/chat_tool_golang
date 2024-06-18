import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

import { environment } from '../../environments/environment';
import { FrontFile } from '../data.service';

@Component({
  selector: 'app-file',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './file.component.html',
  styleUrl: './file.component.scss',
})
export class FileComponent {
  @Input() file: FrontFile = { name: '', size: 0, type: '', ext: '' };

  s3Prefix = environment.S3_PREFIX;
}
