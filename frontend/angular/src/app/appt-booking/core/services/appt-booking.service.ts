import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

// Service interfaces
export interface Service {
  id: number;
  name: string;
  description: string;
  duration_min: number;
  price_cents: number;
  created_at: string;
  updated_at: string;
}

export interface ServiceRequest {
  name: string;
  description: string;
  duration_min: number;
  price_cents: number;
}

export interface ServiceResponse {
  id: number;
  name: string;
  description: string;
  duration_min: number;
  price_cents: number;
  created_at: string;
  updated_at: string;
}

// Staff interfaces
export interface Staff {
  id: number;
  name: string;
  email: string;
  phone: string;
  created_at: string;
  updated_at: string;
}

export interface StaffRequest {
  name: string;
  email: string;
  phone: string;
}

export interface StaffResponse {
  id: number;
  name: string;
  email: string;
  phone: string;
  created_at: string;
  updated_at: string;
}

// Schedule interfaces
export interface Schedule {
  id: number;
  staff_id: number;
  day_of_week: number;
  start_time: string;
  end_time: string;
  created_at: string;
  updated_at: string;
}

export interface ScheduleRequest {
  staff_id: number;
  day_of_week: number;
  start_time: string;
  end_time: string;
}

export interface ScheduleResponse {
  id: number;
  staff_id: number;
  day_of_week: number;
  start_time: string;
  end_time: string;
  created_at: string;
  updated_at: string;
}

// Appointment interfaces
export interface Appointment {
  id: number;
  customer_name: string;
  customer_email: string;
  customer_phone: string;
  staff_id: number;
  service_id: number;
  appointment_datetime: string;
  duration_minutes: number;
  status: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

// AppointmentWithDetails includes joined data from service and staff tables
export interface AppointmentWithDetails extends Appointment {
  service_name: string;
  staff_name: string;
  price_cents: number;
}

export interface BookAppointmentRequest {
  customer_name: string;
  customer_email: string;
  customer_phone: string;
  staff_id: number;
  service_id: number;
  appointment_datetime: string; // ISO 8601 format
  notes?: string;
}

export interface AppointmentResponse {
  id: number;
  customer_name: string;
  customer_email: string;
  customer_phone: string;
  staff_id: number;
  service_id: number;
  appointment_datetime: string;
  duration_minutes: number;
  status: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

@Injectable({
  providedIn: 'root'
})
export class ApptBookingService {
  private baseUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  // ========== SERVICES ==========
  getAllServices(): Observable<ServiceResponse[]> {
    return this.http.get<ServiceResponse[]>(`${this.baseUrl}/api/appt_booking/services`).pipe(
      catchError(this.handleError<ServiceResponse[]>('getAllServices', []))
    );
  }

  getServiceById(id: number): Observable<ServiceResponse> {
    return this.http.get<ServiceResponse>(`${this.baseUrl}/api/appt_booking/services/${id}`).pipe(
      catchError(this.handleError<ServiceResponse>(`getServiceById(${id})`, undefined))
    );
  }

  createService(service: ServiceRequest): Observable<ServiceResponse> {
    return this.http.post<ServiceResponse>(`${this.baseUrl}/api/appt_booking/services`, service).pipe(
      catchError(this.handleError<ServiceResponse>('createService', undefined))
    );
  }

  updateService(id: number, service: ServiceRequest): Observable<ServiceResponse> {
    return this.http.put<ServiceResponse>(`${this.baseUrl}/api/appt_booking/services/${id}`, service).pipe(
      catchError(this.handleError<ServiceResponse>(`updateService(${id})`, undefined))
    );
  }

  deleteService(id: number): Observable<any> {
    return this.http.delete(`${this.baseUrl}/api/appt_booking/services/${id}`).pipe(
      catchError(this.handleError<any>('deleteService', undefined))
    );
  }

  // ========== STAFF ==========
  getAllStaff(): Observable<StaffResponse[]> {
  	return this.http.get<StaffResponse[]>(`${this.baseUrl}/api/appt_booking/staff`).pipe(
  		catchError(this.handleError<StaffResponse[]>('getAllStaff', []))
  	);
  }
 
  getStaffByService(serviceId: number): Observable<StaffResponse[]> {
  	return this.http.get<StaffResponse[]>(`${this.baseUrl}/api/appt_booking/staff/by-service/${serviceId}`).pipe(
  		catchError(this.handleError<StaffResponse[]>(`getStaffByService(${serviceId})`, []))
  	);
  }
 
  getStaffById(id: number): Observable<StaffResponse> {
    return this.http.get<StaffResponse>(`${this.baseUrl}/api/appt_booking/staff/${id}`).pipe(
      catchError(this.handleError<StaffResponse>(`getStaffById(${id})`, undefined))
    );
  }

  createStaff(staff: StaffRequest): Observable<StaffResponse> {
    return this.http.post<StaffResponse>(`${this.baseUrl}/api/appt_booking/staff`, staff).pipe(
      catchError(this.handleError<StaffResponse>('createStaff', undefined))
    );
  }

  updateStaff(id: number, staff: StaffRequest): Observable<StaffResponse> {
    return this.http.put<StaffResponse>(`${this.baseUrl}/api/appt_booking/staff/${id}`, staff).pipe(
      catchError(this.handleError<StaffResponse>(`updateStaff(${id})`, undefined))
    );
  }

  deleteStaff(id: number): Observable<any> {
    return this.http.delete(`${this.baseUrl}/api/appt_booking/staff/${id}`).pipe(
      catchError(this.handleError<any>('deleteStaff', undefined))
    );
  }

  // ========== SCHEDULES ==========
  getAllSchedules(): Observable<ScheduleResponse[]> {
    return this.http.get<ScheduleResponse[]>(`${this.baseUrl}/api/appt_booking/schedules`).pipe(
      catchError(this.handleError<ScheduleResponse[]>('getAllSchedules', []))
    );
  }

  getScheduleById(id: number): Observable<ScheduleResponse> {
    return this.http.get<ScheduleResponse>(`${this.baseUrl}/api/appt_booking/schedules/${id}`).pipe(
      catchError(this.handleError<ScheduleResponse>(`getScheduleById(${id})`, undefined))
    );
  }

  getSchedulesByStaff(staffId: number): Observable<ScheduleResponse[]> {
    return this.http.get<ScheduleResponse[]>(`${this.baseUrl}/api/appt_booking/schedules/staff/${staffId}`).pipe(
      catchError(this.handleError<ScheduleResponse[]>(`getSchedulesByStaff(${staffId})`, []))
    );
  }

  createSchedule(schedule: ScheduleRequest): Observable<ScheduleResponse> {
    return this.http.post<ScheduleResponse>(`${this.baseUrl}/api/appt_booking/schedules`, schedule).pipe(
      catchError(this.handleError<ScheduleResponse>('createSchedule', undefined))
    );
  }

  updateSchedule(id: number, schedule: ScheduleRequest): Observable<ScheduleResponse> {
    return this.http.put<ScheduleResponse>(`${this.baseUrl}/api/appt_booking/schedules/${id}`, schedule).pipe(
      catchError(this.handleError<ScheduleResponse>(`updateSchedule(${id})`, undefined))
    );
  }

  deleteSchedule(id: number): Observable<any> {
    return this.http.delete(`${this.baseUrl}/api/appt_booking/schedules/${id}`).pipe(
      catchError(this.handleError<any>('deleteSchedule', undefined))
    );
  }

  // ========== APPOINTMENTS ==========
   getAllAppointments(staffId?: number, customerEmail?: string): Observable<AppointmentWithDetails[]> {
     let url = `${this.baseUrl}/api/appt_booking/appointments`;
     let params = new HttpParams();
     
     if (staffId) {
       params = params.set('staff_id', staffId.toString());
     }
     if (customerEmail) {
       params = params.set('email', customerEmail);
     }

     return this.http.get<AppointmentWithDetails[]>(url, { params }).pipe(
       catchError(this.handleError<AppointmentWithDetails[]>('getAllAppointments', []))
     );
   }

  getAppointmentById(id: number): Observable<AppointmentResponse> {
    return this.http.get<AppointmentResponse>(`${this.baseUrl}/api/appt_booking/appointments/${id}`).pipe(
      catchError(this.handleError<AppointmentResponse>(`getAppointmentById(${id})`, undefined))
    );
  }

  bookAppointment(appointment: BookAppointmentRequest): Observable<AppointmentResponse> {
    return this.http.post<AppointmentResponse>(`${this.baseUrl}/api/appt_booking/appointments`, appointment).pipe(
      catchError(this.handleError<AppointmentResponse>('bookAppointment', undefined))
    );
  }

  cancelAppointment(id: number): Observable<any> {
    return this.http.put(`${this.baseUrl}/api/appt_booking/appointments/${id}/cancel`, {}).pipe(
      catchError(this.handleError<any>('cancelAppointment', undefined))
    );
  }

  completeAppointment(id: number): Observable<any> {
    return this.http.put(`${this.baseUrl}/api/appt_booking/appointments/${id}/complete`, {}).pipe(
      catchError(this.handleError<any>('completeAppointment', undefined))
    );
  }

  // ========== HELPER METHODS ==========
  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(`${operation} failed:`, error);
      // Let the app keep running by returning a safe result
      return of(result as T);
    };
  }

  // Helper to convert price from dollars to cents
  convertDollarsToCents(dollars: number): number {
    return Math.round(dollars * 100);
  }

  // Helper to convert price from cents to dollars
  convertCentsToDollars(cents: number): number {
    return cents / 100;
  }

  // Helper to format appointment datetime for display
  formatAppointmentDateTime(isoString: string): string {
    const date = new Date(isoString);
    return date.toLocaleString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZoneName: 'short'
    });
  }

  // Helper to format time (HH:MM) for display
  formatTime(time24: string): string {
    const [hours, minutes] = time24.split(':');
    const hour = parseInt(hours, 10);
    const ampm = hour >= 12 ? 'PM' : 'AM';
    const hour12 = hour % 12 || 12;
    return `${hour12}:${minutes} ${ampm}`;
  }

  // Helper to get day of week name
  getDayOfWeekName(day: number): string {
    const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    return days[day] || 'Unknown';
  }

  // Helper to format price in dollars
  formatPrice(cents: number): string {
    return `$${(cents / 100).toFixed(2)}`;
  }

  // Helper to convert local date+time to UTC ISO string for API requests
  convertLocalToUTC(date: string, time: string): string {
    const localDateTime = new Date(`${date}T${time}:00`);
    return localDateTime.toISOString();
  }
 
  // Helper to format local date+time for display (combines conversion + formatting)
  formatLocalAppointment(date: string, time: string): string {
    const utcDateTime = this.convertLocalToUTC(date, time);
    return this.formatAppointmentDateTime(utcDateTime);
  }
}
