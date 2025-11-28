-- Migration: fix_employee_appointments_endpoint_path
-- Created at: 20251128095620
-- Description: Rollback the endpoint path change (not recommended unless reverting code)

-- Revert the endpoint path for GetEmployeeAppointmentsById
UPDATE public.endpoints
SET path = '/employee/:id/appointments'
WHERE method = 'GET' 
  AND path = '/employee/:employee_id/appointments'
  AND controller_name = 'GetEmployeeAppointmentsById';
