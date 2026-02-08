import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ApiService, AboutInfo, HealthStatus } from './api.service';
import { take } from 'rxjs/operators';

describe('ApiService', () => {
  let service: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [ApiService]
    });
    service = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should return healthy true when health endpoint returns ok', (done) => {
    const mockResponse = { status: 'ok' };
    const expectedResult: HealthStatus = { healthy: true };

    service.checkHealth().pipe(take(1)).subscribe(result => {
      expect(result).toEqual(expectedResult);
      done();
    });

    const req = httpMock.expectOne(`${service['baseUrl']}/health`);
    expect(req.request.method).toBe('GET');
    req.flush(mockResponse);
  });

  it('should return healthy false when health endpoint returns non-ok', (done) => {
    const mockResponse = { status: 'error' };
    const expectedResult: HealthStatus = { healthy: false };

    service.checkHealth().pipe(take(1)).subscribe(result => {
      expect(result).toEqual(expectedResult);
      done();
    });

    const req = httpMock.expectOne(`${service['baseUrl']}/health`);
    expect(req.request.method).toBe('GET');
    req.flush(mockResponse);
  });

  it('should return healthy false on network error', (done) => {
    service.checkHealth().pipe(take(1)).subscribe(result => {
      expect(result).toEqual({ healthy: false });
      done();
    });

    const req = httpMock.expectOne(`${service['baseUrl']}/health`);
    expect(req.request.method).toBe('GET');
    req.error(new Error('Network error'));
  });

  it('should return info data on successful fetch', (done) => {
    const mockResponse = { version: '1.0.0', uptime: 123 };
    const expectedResult = mockResponse;

    service.getInfo().pipe(take(1)).subscribe(result => {
      expect(result).toEqual(expectedResult);
      done();
    });

    const req = httpMock.expectOne(`${service['baseUrl']}/info`);
    expect(req.request.method).toBe('GET');
    req.flush(mockResponse);
  });

  it('should return null for info on fetch error', (done) => {
    service.getInfo().pipe(take(1)).subscribe(result => {
      expect(result).toBeNull();
      done();
    });

    const req = httpMock.expectOne(`${service['baseUrl']}/info`);
    expect(req.request.method).toBe('GET');
    req.error(new Error('Network error'));
  });

  it('should return about data', (done) => {
    const expectedResult: AboutInfo = {
      status: 'healthy',
      message: 'K8s Fullstack Blueprint'
    };

    service.getAbout().pipe(take(1)).subscribe(result => {
      expect(result).toEqual(expectedResult);
      done();
    });

    // No HTTP request expected for getAbout() since it returns static data
    httpMock.verify();
  });
});
