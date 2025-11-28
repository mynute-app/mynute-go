-- Migration: change_employee_endpoint_path_parameters
-- Created at: 20251128111531
-- Description: Rollback the endpoint path changes (not recommended unless reverting code)

-- Employee Endpoints: Revert :employee_id back to :id
UPDATE public.endpoints
SET path = '/employee/:id/work_schedule'
WHERE method = 'POST' 
  AND path = '/employee/:employee_id/work_schedule'
  AND controller_name = 'CreateEmployeeWorkSchedule';

UPDATE public.endpoints
SET path = '/employee/:id/work_schedule'
WHERE method = 'GET' 
  AND path = '/employee/:employee_id/work_schedule'
  AND controller_name = 'GetEmployeeWorkSchedule';

UPDATE public.endpoints
SET path = '/employee/:id/work_range/:work_range_id'
WHERE method = 'GET' 
  AND path = '/employee/:employee_id/work_range/:work_range_id'
  AND controller_name = 'GetEmployeeWorkRange';

UPDATE public.endpoints
SET path = '/employee/:id/work_range/:work_range_id'
WHERE method = 'DELETE' 
  AND path = '/employee/:employee_id/work_range/:work_range_id'
  AND controller_name = 'DeleteEmployeeWorkRange';

UPDATE public.endpoints
SET path = '/employee/:id/work_range/:work_range_id'
WHERE method = 'PUT' 
  AND path = '/employee/:employee_id/work_range/:work_range_id'
  AND controller_name = 'UpdateEmployeeWorkRange';

UPDATE public.endpoints
SET path = '/employee/:id/work_range/:work_range_id/services'
WHERE method = 'POST' 
  AND path = '/employee/:employee_id/work_range/:work_range_id/services'
  AND controller_name = 'AddEmployeeWorkRangeServices';

UPDATE public.endpoints
SET path = '/employee/:id/work_range/:work_range_id/service/:service_id'
WHERE method = 'DELETE' 
  AND path = '/employee/:employee_id/work_range/:work_range_id/service/:service_id'
  AND controller_name = 'DeleteEmployeeWorkRangeService';

UPDATE public.endpoints
SET path = '/employee/:id/appointments'
WHERE method = 'GET' 
  AND path = '/employee/:employee_id/appointments'
  AND controller_name = 'GetEmployeeAppointmentsById';
