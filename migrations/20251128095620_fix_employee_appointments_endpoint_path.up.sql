-- Migration: fix_employee_appointments_endpoint_path
-- Created at: 20251128095620
-- Description: Updates the employee appointments endpoint path from :id to :employee_id
--              to match the current code definition and ensure authorization works correctly

-- Update the endpoint path for GetEmployeeAppointmentsById
UPDATE public.endpoints
SET path = '/employee/:employee_id/appointments'
WHERE method = 'GET' 
  AND path = '/employee/:id/appointments'
  AND controller_name = 'GetEmployeeAppointmentsById';
