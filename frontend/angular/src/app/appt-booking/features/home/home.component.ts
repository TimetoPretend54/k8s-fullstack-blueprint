import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-appt-booking-home',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="home-container">
      <h1>Welcome to Appointment Booking</h1>
      <p>This is a proof-of-concept (POC) appointment booking system.</p>
      
      <div class="feature-cards">
        <div class="card">
          <h3>Book an Appointment</h3>
          <p>Schedule a new appointment with our staff</p>
          <a routerLink="/appt-booking/booking" class="btn">Start Booking</a>
        </div>
        
        <div class="card">
          <h3>My Appointments</h3>
          <p>View and manage your upcoming appointments</p>
          <a routerLink="/appt-booking/my-appointments" class="btn">View Appointments</a>
        </div>
        
        <div class="card">
          <h3>Dashboard</h3>
          <p>Overview of appointments and services</p>
          <a routerLink="/appt-booking/dashboard" class="btn">Go to Dashboard</a>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .home-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 40px 20px;
      text-align: center;
    }
    
    h1 {
      color: #1976d2;
      margin-bottom: 10px;
    }
    
    p {
      color: #666;
      margin-bottom: 40px;
    }
    
    .feature-cards {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 30px;
      margin-top: 20px;
    }
    
    .card {
      background: white;
      border-radius: 8px;
      padding: 30px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
      transition: transform 0.2s, box-shadow 0.2s;
    }
    
    .card:hover {
      transform: translateY(-5px);
      box-shadow: 0 4px 16px rgba(0,0,0,0.15);
    }
    
    .card h3 {
      color: #333;
      margin-bottom: 15px;
    }
    
    .card p {
      color: #666;
      margin-bottom: 20px;
      font-size: 14px;
    }
    
    .btn {
      display: inline-block;
      background-color: #1976d2;
      color: white;
      padding: 12px 24px;
      border-radius: 4px;
      text-decoration: none;
      transition: background-color 0.2s;
    }
    
    .btn:hover {
      background-color: #1565c0;
    }
  `]
})
export class ApptBookingHomeComponent {}
