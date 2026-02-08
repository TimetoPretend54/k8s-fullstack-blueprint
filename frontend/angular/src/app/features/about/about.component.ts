import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from '../../core/services/api.service';
import { AboutInfo } from '../../core/services/api.service';

@Component({
  selector: 'app-about',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class AboutComponent implements OnInit {
  aboutData?: AboutInfo;

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.apiService.getAbout().subscribe(data => {
      this.aboutData = data;
    });
  }
}
