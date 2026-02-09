import { Routes } from '@angular/router';
import { ApptBookingHomeComponent } from './features/home/home.component';
import { BookingComponent } from './features/booking/booking.component';
import { MyAppointmentsComponent } from './features/my-appointments/my-appointments.component';
import { DashboardComponent } from './features/dashboard/dashboard.component';

export const apptBookingRoutes: Routes = [
  {
    path: '',
    redirectTo: 'home',
    pathMatch: 'full'
  },
  {
    path: 'home',
    component: ApptBookingHomeComponent
  },
  {
    path: 'booking',
    component: BookingComponent
  },
  {
    path: 'my-appointments',
    component: MyAppointmentsComponent
  },
  {
    path: 'dashboard',
    component: DashboardComponent
  }
];
