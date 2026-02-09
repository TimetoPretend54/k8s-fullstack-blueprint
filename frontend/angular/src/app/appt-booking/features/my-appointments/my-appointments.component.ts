import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApptBookingService } from '../../core/services/appt-booking.service';
import { AppointmentWithDetails } from '../../core/services/appt-booking.service';

@Component({
  selector: 'app-my-appointments',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './my-appointments.component.html',
  styleUrls: ['./my-appointments.component.scss']
})
export class MyAppointmentsComponent implements OnInit {
   // Form state
   email: string = '';
   searchPerformed: boolean = false;
   
   // Data
   appointments: AppointmentWithDetails[] = [];
   loading: boolean = false;
   error: string | null = null;
   
   // Filtering
   upcomingOnly: boolean = true;
   
   // Email suggestions
   availableEmails: string[] = [];
   
   constructor(private apptBookingService: ApptBookingService) {}
   
   ngOnInit(): void {
     // Load all appointments to extract unique emails for suggestions
     this.loadEmailSuggestions();
   }
   
   loadEmailSuggestions(): void {
     this.apptBookingService.getAllAppointments()
       .subscribe({
         next: (appts: AppointmentWithDetails[]) => {
           // Extract unique customer emails
           const emailSet = new Set<string>();
           appts.forEach(appt => {
             if (appt.customer_email) {
               emailSet.add(appt.customer_email);
             }
           });
           this.availableEmails = Array.from(emailSet).sort();
         },
         error: (err) => {
           console.error('Error loading email suggestions:', err);
           // Non-critical - just won't show suggestions
         }
       });
   }
  
  searchAppointments(): void {
    if (!this.email) {
      this.error = 'Please enter your email address';
      return;
    }
    
    this.loading = true;
    this.error = null;
    this.searchPerformed = true;
    
    this.apptBookingService.getAllAppointments(undefined, this.email)
      .subscribe({
        next: (appts: AppointmentWithDetails[]) => {
          this.appointments = appts;
          this.loading = false;
        },
        error: (err: any) => {
          this.error = 'Failed to load appointments. Please try again.';
          this.loading = false;
          console.error('Error loading appointments:', err);
        }
      });
  }
  
  cancelAppointment(appointmentId: number): void {
    if (!confirm('Are you sure you want to cancel this appointment?')) {
      return;
    }
    
    this.loading = true;
    this.apptBookingService.cancelAppointment(appointmentId)
      .subscribe({
        next: () => {
          // Remove the cancelled appointment from the list
          this.appointments = this.appointments.filter(a => a.id !== appointmentId);
          this.loading = false;
        },
        error: (err) => {
          this.error = 'Failed to cancel appointment. Please try again.';
          this.loading = false;
          console.error('Error cancelling appointment:', err);
        }
      });
  }
  
  getStatusBadgeClass(status: string): string {
    switch (status) {
      case 'confirmed':
        return 'badge-confirmed';
      case 'cancelled':
        return 'badge-cancelled';
      case 'completed':
        return 'badge-completed';
      default:
        return '';
    }
  }
  
  formatDateTime(dateTimeStr: string): string {
    const date = new Date(dateTimeStr);
    return date.toLocaleString('en-US', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  
  formatPrice(priceCents: number): string {
    return this.apptBookingService.formatPrice(priceCents);
  }
  
  getUpcomingAppointments(): AppointmentWithDetails[] {
    if (!this.upcomingOnly) {
      return this.appointments;
    }
    const now = new Date();
    return this.appointments.filter(a => 
      a.status === 'confirmed' && new Date(a.appointment_datetime) > now
    );
  }
   
  getPastAppointments(): AppointmentWithDetails[] {
    return this.appointments.filter(a => 
      a.status === 'completed' || a.status === 'cancelled' || new Date(a.appointment_datetime) <= new Date()
    );
  }
}
