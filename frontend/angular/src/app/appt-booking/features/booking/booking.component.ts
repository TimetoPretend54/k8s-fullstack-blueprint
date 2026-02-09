import { Component, OnInit } from '@angular/core';
import { CommonModule, ViewportScroller } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApptBookingService, Service as ServiceModel, Staff, Schedule, Appointment } from '../../core/services/appt-booking.service';
import { Observable } from 'rxjs';

export type BookingStep = 'service' | 'staff' | 'schedule' | 'confirmation';

@Component({
  selector: 'app-booking',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './booking.component.html',
  styleUrls: ['./booking.component.scss']
})
export class BookingComponent implements OnInit {
  currentStep: BookingStep = 'service';
  
  // Data from API
  services$!: Observable<ServiceModel[]>;
  staff$!: Observable<Staff[]>;
  schedules$!: Observable<Schedule[]>;
  
  // Selections
  selectedService: ServiceModel | null = null;
  selectedStaff: Staff | null = null;
  selectedDate: string = ''; // YYYY-MM-DD format
  selectedTime: string = ''; // HH:MM format
  selectedSchedule: Schedule | null = null;
  
  // Customer info
  customerName: string = '';
  customerEmail: string = '';
  customerPhone: string = '';
  
  // UI state
  isLoading: boolean = false;
  errorMessage: string | null = null;
  successMessage: string | null = null;
  
   constructor(public apptBookingService: ApptBookingService, private viewportScroller: ViewportScroller) {
     this.loadServices();
     this.loadStaff();
   }
  
  ngOnInit(): void {
    // Services and staff are loaded in constructor
  }
  
  loadServices(): void {
    this.services$ = this.apptBookingService.getAllServices();
  }
  
  loadStaff(serviceId?: number): void {
  	if (serviceId) {
  		this.staff$ = this.apptBookingService.getStaffByService(serviceId);
  	} else {
  		this.staff$ = this.apptBookingService.getAllStaff();
  	}
  }
  
  loadSchedules(): void {
  	if (this.selectedStaff) {
  		this.schedules$ = this.apptBookingService.getSchedulesByStaff(this.selectedStaff.id);
  	} else {
  		this.schedules$ = new Observable<Schedule[]>(observer => observer.next([]));
  	}
  }
  
  getSchedulesForDate(schedules: Schedule[], date: string): Schedule[] {
  	if (!date || !schedules) return [];
  	const selectedDate = new Date(date + 'T00:00:00');
  	const selectedDayOfWeek = selectedDate.getDay(); // 0 (Sunday) to 6 (Saturday)
  	return schedules.filter(schedule => schedule.day_of_week === selectedDayOfWeek);
  }
  
   /** Generate time slots within a schedule based on service duration */
  	getTimeSlots(schedule: Schedule): string[] {
   	console.log('getTimeSlots called', { schedule, selectedService: this.selectedService });
   	if (!schedule || !this.selectedService) {
   		console.log('  returning [] due to missing schedule or selectedService');
   		return [];
   	}
   	
   	// Parse start and end times (format: "HH:MM")
   	const [startHour, startMin] = schedule.start_time.split(':').map(Number);
   	const [endHour, endMin] = schedule.end_time.split(':').map(Number);
   	
   	const startMinutes = startHour * 60 + startMin;
   	const endMinutes = endHour * 60 + endMin;
   	const duration = this.selectedService.duration_min;
   	console.log(`  Parsed: start=${startMinutes}min, end=${endMinutes}min, duration=${duration}min`);
   	
   	const timeSlots: string[] = [];
   	for (let current = startMinutes; current + duration <= endMinutes; current += duration) {
   		const hour = Math.floor(current / 60);
   		const minute = current % 60;
   		const timeString = `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`;
   		timeSlots.push(timeString);
   	}
   	console.log('  returning', timeSlots);
   	return timeSlots;
   }
  
  selectService(service: ServiceModel): void {
  	this.selectedService = service;
  	this.selectedStaff = null;
  	this.selectedSchedule = null;
  	this.selectedDate = '';
  	this.selectedTime = '';
  	this.loadStaff(service.id);
  	this.currentStep = 'staff';
  }
  
  selectStaff(staff: Staff): void {
    this.selectedStaff = staff;
    this.selectedSchedule = null;
    this.selectedDate = '';
    this.selectedTime = '';
    this.loadSchedules();
    this.currentStep = 'schedule';
  }
  
  selectDate(date: string): void {
    this.selectedDate = date;
    this.selectedTime = '';
    this.selectedSchedule = null;
  }
  
  selectTime(time: string): void {
    this.selectedTime = time;
    this.findOrCreateSchedule();
  }
  
   findOrCreateSchedule(): void {
   	if (!this.selectedDate || !this.selectedTime || !this.selectedStaff || !this.selectedService) {
   		return;
   	}
   	
   	// Get day of week from selected date (0 = Sunday, 6 = Saturday)
   	const selectedDateObj = new Date(this.selectedDate + 'T00:00:00');
   	const selectedDayOfWeek = selectedDateObj.getDay();
   	
   	// Find matching schedule: same day of week and time falls within schedule's start-end
   	this.schedules$.subscribe(schedules => {
   		const matchingSchedule = schedules.find(s => {
   			return s.day_of_week === selectedDayOfWeek &&
   			       this.selectedTime! >= s.start_time &&
   			       this.selectedTime! < s.end_time;
   		});
   		
   		if (matchingSchedule) {
   			this.selectedSchedule = matchingSchedule;
   		} else {
   			// No schedule found - user needs to select a different time
   			this.selectedSchedule = null;
   			this.errorMessage = 'No schedule available for this time. Please select a different time.';
   		}
   	});
   }
  
   canProceedToConfirmation(): boolean {
     const hasSelections = this.selectedService !== null && 
                           this.selectedStaff !== null && 
                           this.selectedSchedule !== null && 
                           this.selectedDate !== '' && 
                           this.selectedTime !== '';
     const hasCustomerInfo = this.customerName.trim() !== '' && 
                             this.customerEmail.trim() !== '' && 
                             this.customerPhone.trim() !== '';
     return hasSelections && hasCustomerInfo;
   }
   
   canProceedFromSchedule(): boolean {
     return this.selectedService !== null && 
            this.selectedStaff !== null && 
            this.selectedSchedule !== null && 
            this.selectedDate !== '' && 
            this.selectedTime !== '';
   }
  
  getTotalPrice(): number {
    return this.selectedService ? this.selectedService.price_cents : 0;
  }
  
   getAppointmentDateTime(): string {
     if (!this.selectedDate || !this.selectedTime) return '';
     // Use service method to convert local time to UTC
     return this.apptBookingService.convertLocalToUTC(this.selectedDate, this.selectedTime);
   }
  
  bookAppointment(): void {
    if (!this.canProceedToConfirmation()) {
      this.errorMessage = 'Please fill in all required fields.';
      return;
    }
    
    this.isLoading = true;
    this.errorMessage = null;
    
    const appointmentData = {
      customer_name: this.customerName.trim(),
      customer_email: this.customerEmail.trim(),
      customer_phone: this.customerPhone.trim(),
      staff_id: this.selectedStaff!.id,
      service_id: this.selectedService!.id,
      appointment_datetime: this.getAppointmentDateTime()
    };
    
     this.apptBookingService.bookAppointment(appointmentData).subscribe({
       next: (appointment) => {
         this.isLoading = false;
         this.successMessage = 'Appointment booked successfully!';
         this.currentStep = 'confirmation';
         // Scroll to top after successful booking
         this.viewportScroller.scrollToPosition([0, 0]);
       },
       error: (err) => {
         this.isLoading = false;
         this.errorMessage = err.error?.message || 'Failed to book appointment. Please try again.';
         console.error(err);
       }
     });
  }
  
  goToStep(step: BookingStep): void {
    this.currentStep = step;
    this.errorMessage = null;
  }
  
  resetBooking(): void {
    this.currentStep = 'service';
    this.selectedService = null;
    this.selectedStaff = null;
    this.selectedSchedule = null;
    this.selectedDate = '';
    this.selectedTime = '';
    this.customerName = '';
    this.customerEmail = '';
    this.customerPhone = '';
    this.errorMessage = null;
    this.successMessage = null;
  }
}
