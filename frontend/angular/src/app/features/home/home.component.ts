import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService, DemoData } from '../../core/services/api.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  demoDataList: DemoData[] = [];
  newContent: string = '';
  isLoading: boolean = false;
  errorMessage: string | null = null;

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.loadDemoData();
  }

  loadDemoData(): void {
    this.isLoading = true;
    this.errorMessage = null;
    
    this.apiService.getAllDemoData().subscribe({
      next: (data) => {
        this.demoDataList = data;
        this.isLoading = false;
      },
      error: (err) => {
        this.errorMessage = 'Failed to load demo data';
        this.isLoading = false;
        console.error(err);
      }
    });
  }

  onSubmit(): void {
    if (!this.newContent.trim()) {
      return;
    }

    this.isLoading = true;
    this.errorMessage = null;

    this.apiService.upsertDemoData(this.newContent).subscribe({
      next: () => {
        this.newContent = '';
        this.loadDemoData(); // Refresh the list
      },
      error: (err) => {
        this.errorMessage = 'Failed to save demo data';
        this.isLoading = false;
        console.error(err);
      }
    });
  }
}
