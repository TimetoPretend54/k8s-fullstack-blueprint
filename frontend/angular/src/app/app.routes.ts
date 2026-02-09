import { Routes } from '@angular/router';
import { apptBookingRoutes } from './appt-booking/appt-booking.routes';

export const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./features/home/home.routes').then(m => m.homeRoutes)
  },
  {
    path: 'about',
    loadChildren: () => import('./features/about/about.routes').then(m => m.aboutRoutes)
  },
  {
    path: 'appt-booking',
    children: apptBookingRoutes
  }
];
