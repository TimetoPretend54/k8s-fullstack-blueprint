import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  template: `
    <nav>
      <div class="nav-section">
        <span class="section-label">Main Application</span>
        <a routerLink="/" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">Home</a>
        <a routerLink="/about" routerLinkActive="active">About</a>
      </div>
      <span class="divider"></span>
      <div class="nav-section">
        <span class="section-label">Appointment Booking</span>
        <a routerLink="/appt-booking/home" routerLinkActive="active">Booking Home</a>
        <a routerLink="/appt-booking/booking" routerLinkActive="active">Book Appointment</a>
        <a routerLink="/appt-booking/my-appointments" routerLinkActive="active">My Appointments</a>
        <a routerLink="/appt-booking/dashboard" routerLinkActive="active">Dashboard</a>
      </div>
    </nav>
  `,
  styles: [`
    nav {
      display: flex;
      gap: 30px;
      justify-content: center;
      align-items: center;
      flex-wrap: wrap;
      padding: 1rem 0;
    }

    .nav-section {
      display: flex;
      gap: 15px;
      align-items: center;
      flex-wrap: wrap;
    }

    .section-label {
      color: rgba(255, 255, 255, 0.7);
      font-size: 0.85rem;
      font-weight: 500;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      padding-right: 8px;
    }

    nav a {
      color: white;
      text-decoration: none;
      padding: 8px 16px;
      border-radius: 4px;
      transition: background-color 0.2s;
      font-weight: 500;
    }

    nav a:hover {
      background-color: rgba(255, 255, 255, 0.1);
    }

    nav a.active {
      background-color: rgba(255, 255, 255, 0.2);
      color: white;
    }

    .divider {
      width: 1px;
      height: 30px;
      background-color: rgba(255, 255, 255, 0.3);
    }
  `]
})
export class NavbarComponent {}
