USE gate_keeper_demo;

CREATE TABLE IF NOT EXISTS application_gates (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    cycle_year INT NOT NULL,
    gate_code VARCHAR(64) NOT NULL,
    sequence_number INT NOT NULL,
    is_active CHAR(1) NOT NULL DEFAULT 'Y',
    open_date DATETIME NOT NULL,
    close_date DATETIME NOT NULL,
    notes VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_application_gates_year_gate (cycle_year, gate_code),
    UNIQUE KEY uq_application_gates_year_sequence (cycle_year, sequence_number),
    KEY idx_application_gates_active_year (is_active, cycle_year)
);

INSERT INTO application_gates (
    cycle_year,
    gate_code,
    sequence_number,
    is_active,
    open_date,
    close_date,
    notes
) VALUES
    (2025, 'A', 1, 'Y', DATE_ADD(NOW(), INTERVAL -420 DAY), DATE_ADD(NOW(), INTERVAL -418 DAY), 'Past cycle dummy data'),
    (2025, 'B', 2, 'Y', DATE_ADD(NOW(), INTERVAL -416 DAY), DATE_ADD(NOW(), INTERVAL -414 DAY), 'Past cycle dummy data'),
    (2025, 'C', 3, 'Y', DATE_ADD(NOW(), INTERVAL -412 DAY), DATE_ADD(NOW(), INTERVAL -410 DAY), 'Past cycle dummy data'),
    (2025, 'D', 4, 'Y', DATE_ADD(NOW(), INTERVAL -408 DAY), DATE_ADD(NOW(), INTERVAL -406 DAY), 'Past cycle dummy data'),
    (2025, 'E', 5, 'Y', DATE_ADD(NOW(), INTERVAL -404 DAY), DATE_ADD(NOW(), INTERVAL -402 DAY), 'Past cycle dummy data'),
    (2026, 'A', 1, 'Y', DATE_ADD(NOW(), INTERVAL -20 DAY), DATE_ADD(NOW(), INTERVAL -18 DAY), 'Current cycle dummy data'),
    (2026, 'B', 2, 'Y', DATE_ADD(NOW(), INTERVAL -16 DAY), DATE_ADD(NOW(), INTERVAL -14 DAY), 'Current cycle dummy data'),
    (2026, 'C', 3, 'Y', DATE_ADD(NOW(), INTERVAL -12 DAY), DATE_ADD(NOW(), INTERVAL -10 DAY), 'Current cycle dummy data'),
    (2026, 'D', 4, 'Y', DATE_ADD(NOW(), INTERVAL -1 DAY), DATE_ADD(NOW(), INTERVAL 1 DAY), 'Current cycle dummy data'),
    (2026, 'E', 5, 'Y', DATE_ADD(NOW(), INTERVAL 4 DAY), DATE_ADD(NOW(), INTERVAL 6 DAY), 'Current cycle dummy data'),
    (2026, 'F', 6, 'N', DATE_ADD(NOW(), INTERVAL 8 DAY), DATE_ADD(NOW(), INTERVAL 10 DAY), 'Inactive gate to demonstrate filtering'),
    (2027, 'A', 1, 'Y', DATE_ADD(NOW(), INTERVAL 120 DAY), DATE_ADD(NOW(), INTERVAL 122 DAY), 'Future cycle dummy data'),
    (2027, 'B', 2, 'Y', DATE_ADD(NOW(), INTERVAL 124 DAY), DATE_ADD(NOW(), INTERVAL 126 DAY), 'Future cycle dummy data'),
    (2027, 'C', 3, 'Y', DATE_ADD(NOW(), INTERVAL 128 DAY), DATE_ADD(NOW(), INTERVAL 130 DAY), 'Future cycle dummy data'),
    (2027, 'D', 4, 'Y', DATE_ADD(NOW(), INTERVAL 132 DAY), DATE_ADD(NOW(), INTERVAL 134 DAY), 'Future cycle dummy data'),
    (2027, 'E', 5, 'Y', DATE_ADD(NOW(), INTERVAL 136 DAY), DATE_ADD(NOW(), INTERVAL 138 DAY), 'Future cycle dummy data');
