-- +goose Up
SELECT 'up SQL query';
CREATE TYPE appointment_status AS ENUM ('scheduled', 'completed', 'canceled');


CREATE TYPE doctor_status AS ENUM ('ACTIVE', 'ON_LEAVE', 'INACTIVE');

CREATE TABLE roles (
                       id BIGSERIAL PRIMARY KEY,
                       name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE auth (
                      id BIGSERIAL PRIMARY KEY,
                      email VARCHAR(255) UNIQUE NOT NULL,
                      password_hash VARCHAR(255) NOT NULL,
                      active BOOLEAN DEFAULT TRUE,
                      is_seller BOOLEAN DEFAULT FALSE,
                      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                      updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                      deleted_at TIMESTAMP
);


CREATE TABLE auth_roles (
                            auth_id BIGINT NOT NULL REFERENCES auth(id),
                            role_id BIGINT NOT NULL REFERENCES roles(id),
                            PRIMARY KEY (auth_id, role_id)
);


CREATE TABLE sessions (
                          session_key VARCHAR(255) PRIMARY KEY,
                          last_used TIMESTAMP NOT NULL,
                          expires_at TIMESTAMP NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          user_id BIGINT NOT NULL REFERENCES auth(id)
);


CREATE TABLE specializations (
                                 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                 name VARCHAR(255) UNIQUE NOT NULL,
                                 description TEXT
);


CREATE TABLE doctors (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         first_name VARCHAR(255) NOT NULL,
                         last_name VARCHAR(255) NOT NULL,
                         patronymic VARCHAR(255),
                         date_of_birth DATE NOT NULL,
                         specialization_id UUID NOT NULL REFERENCES specializations(id),
                         status doctor_status DEFAULT 'ACTIVE',
                         created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE patients (
                          id BIGSERIAL PRIMARY KEY,
                          user_id UUID NOT NULL,
                          blood_type VARCHAR(5),
                          height DECIMAL(5,2),
                          weight DECIMAL(5,2),
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE allergies (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           patient_id BIGINT NOT NULL REFERENCES patients(id),
                           name VARCHAR(255) NOT NULL,
                           severity VARCHAR(50),
                           created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE insurances (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            patient_id BIGINT NOT NULL REFERENCES patients(id),
                            provider VARCHAR(255) NOT NULL,
                            policy_number VARCHAR(255) UNIQUE NOT NULL,
                            expiration_date DATE NOT NULL
);


CREATE TABLE schedule_slots (
                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                doctor_id UUID NOT NULL REFERENCES doctors(id),
                                date DATE NOT NULL,
                                start_time TIME NOT NULL,
                                end_time TIME NOT NULL,
                                is_available BOOLEAN DEFAULT TRUE,
                                appointment_id UUID,
                                meeting_link TEXT
);


CREATE TABLE appointments (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              user_id BIGINT NOT NULL,
                              doctor_id UUID NOT NULL REFERENCES doctors(id),
                              date TIMESTAMP NOT NULL,
                              status appointment_status NOT NULL DEFAULT 'scheduled',
                              notes TEXT,
                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              deleted_at TIMESTAMP
);


CREATE TABLE prescriptions (
                               id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                               patient_id UUID NOT NULL,
                               doctor_id UUID NOT NULL REFERENCES doctors(id),
                               medication VARCHAR(255) NOT NULL,
                               dosage VARCHAR(255) NOT NULL,
                               valid_until DATE NOT NULL
);


CREATE INDEX idx_appointments_user_id ON appointments(user_id);
CREATE INDEX idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX idx_appointments_date ON appointments(date);
CREATE INDEX idx_schedule_slots_doctor_id_date ON schedule_slots(doctor_id, date);
CREATE INDEX idx_prescriptions_patient_id ON prescriptions(patient_id);

-- +goose Down
-- +goose StatementBegin

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

DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS doctor_status;
-- +goose StatementEnd

