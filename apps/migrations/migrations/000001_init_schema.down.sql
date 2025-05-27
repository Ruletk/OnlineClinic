-- Удаляем таблицы в обратном порядке от создания
DROP TABLE IF EXISTS prescriptions;
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS schedule_slots;
DROP TABLE IF EXISTS insurances;
DROP TABLE IF EXISTS allergies;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS doctors;
DROP TABLE IF EXISTS specializations;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS auth_roles;
DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS roles;

-- Удаляем типы ENUM
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS doctor_status;
