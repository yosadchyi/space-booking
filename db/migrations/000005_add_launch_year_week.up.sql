ALTER TABLE launch ADD week INT NOT NULL DEFAULT '0';
ALTER TABLE launch ADD year INT NOT NULL DEFAULT '0';
CREATE INDEX launch_launchpad_year_week_idx ON launch(launchpad_id, year, week);
