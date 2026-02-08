import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  template: `
    <nav>
      <a routerLink="/" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">Home</a>
      <a routerLink="/about" routerLinkActive="active">About</a>
    </nav>
  `,
  styles: [`
    nav {
      display: flex;
      gap: 20px;
      justify-content: center;
    }

    nav a {
      color: #1976d2;
      text-decoration: none;
      padding: 8px 16px;
      border-radius: 4px;
      transition: background-color 0.2s;
    }

    nav a:hover {
      background-color: #e3f2fd;
    }

    nav a.active {
      background-color: #1976d2;
      color: white;
    }
  `]
})
export class NavbarComponent {}
