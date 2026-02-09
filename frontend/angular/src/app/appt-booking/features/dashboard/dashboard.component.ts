import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApptBookingService } from '../../core/services/appt-booking.service';
import { AppointmentWithDetails, Service, Staff } from '../../core/services/appt-booking.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  // Data
  appointments: AppointmentWithDetails[] = [];
  services: Service[] = [];
  staff: Staff[] = [];
  
  // Loading state
  loading: boolean = false;
  error: string | null = null;
  
  // Stats
  totalRevenue: number = 0;
  appointmentsByStaff: Map<number, number> = new Map();
  
  constructor(private apptBookingService: ApptBookingService) {}
  
  ngOnInit(): void {
    this.loadAllData();
  }
  
  loadAllData(): void {
    this.loading = true;
    this.error = null;
    
    // Load all appointments
    this.apptBookingService.getAllAppointments().subscribe({
      next: (appts) => {
        this.appointments = appts.sort((a, b) => 
          new Date(b.appointment_datetime).getTime() - new Date(a.appointment_datetime).getTime()
        );
        this.calculateStats();
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load dashboard data.';
        this.loading = false;
        console.error('Error loading dashboard:', err);
      }
    });
    
    // Load services and staff for reference
    this.apptBookingService.getAllServices().subscribe({
      next: (svcs) => this.services = svcs,
      error: (err) => console.error('Error loading services:', err)
    });
    
    this.apptBookingService.getAllStaff().subscribe({
      next: (staffList) => this.staff = staffList,
      error: (err) => console.error('Error loading staff:', err)
    });
  }
  
  calculateStats(): void {
    // Calculate total revenue from completed appointments
    this.totalRevenue = 0;
    this.appointmentsByStaff.clear();
    
    for (const appt of this.appointments) {
      if (appt.status === 'completed') {
        this.totalRevenue += appt.price_cents;
      }
      
      // Count appointments by staff
      const currentCount = this.appointmentsByStaff.get(appt.staff_id) || 0;
      this.appointmentsByStaff.set(appt.staff_id, currentCount + 1);
    }
  }
  
  getStaffName(staffId: number): string {
    const staffMember = this.staff.find(s => s.id === staffId);
    return staffMember ? staffMember.name : 'Unknown';
  }
  
  getServiceName(serviceId: number): string {
    const service = this.services.find(s => s.id === serviceId);
    return service ? service.name : 'Unknown';
  }
  
  formatPrice(cents: number): string {
    return this.apptBookingService.formatPrice(cents);
  }
  
  formatDateTime(dateTimeStr: string): string {
    const date = new Date(dateTimeStr);
    return date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
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
  
  cancelAppointment(appointmentId: number): void {
    if (!confirm('Cancel this appointment?')) {
      return;
    }
    
    this.apptBookingService.cancelAppointment(appointmentId).subscribe({
      next: () => {
        this.appointments = this.appointments.filter(a => a.id !== appointmentId);
        this.calculateStats();
      },
      error: (err) => {
        console.error('Error cancelling appointment:', err);
        alert('Failed to cancel appointment.');
      }
    });
  }
  
  completeAppointment(appointmentId: number): void {
    this.apptBookingService.completeAppointment(appointmentId).subscribe({
      next: () => {
        // Update the appointment status locally
        const appt = this.appointments.find(a => a.id === appointmentId);
        if (appt) {
          appt.status = 'completed';
          this.calculateStats();
        }
      },
      error: (err) => {
        console.error('Error completing appointment:', err);
        alert('Failed to mark appointment as complete.');
      }
    });
  }
}
