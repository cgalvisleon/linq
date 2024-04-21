package lib

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// DDL Data Definition Language
// This package contains the functions to definition data elements in the database

// ddlListen return sql used listen ddl
func ddlListen() string {
	return `
	CREATE SCHEMA IF NOT EXISTS linq;

	CREATE OR REPLACE FUNCTION linq.SYNC_INSERT()
  RETURNS
    TRIGGER AS $$
  DECLARE
   CHANNEL VARCHAR(250);
  BEGIN
    IF NEW._IDT = '-1' THEN
      NEW._IDT = uuid_generate_v4();
		END IF;

		CHANNEL = TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
		PERFORM pg_notify(
		CHANNEL,
		json_build_object(
			'option', TG_OP,
			'_idt', NEW._IDT
		)::text
		);
    
  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION linq.SYNC_UPDATE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    CHANNEL VARCHAR(250);
  BEGIN
    IF NEW._IDT = '-1' THEN
			NEW._IDT = uuid_generate_v4();    
    END IF;
    
		CHANNEL = TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
		PERFORM pg_notify(
		CHANNEL,
		json_build_object(
			'option', TG_OP,
			'_idt', NEW._IDT  
		)::text
		);

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE linq.FUNCTION SYNC_DELETE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    CHANNEL VARCHAR(250);
  BEGIN
		CHANNEL = TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
		PERFORM pg_notify(
		CHANNEL,
		json_build_object(
			'option', TG_OP,
			'_idt', OLD._IDT
		)::text
		);

  RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;
	`
}

// ddlSync return sql sync ddl
func ddlSync() string {
	return `
  -- DROP TABLE IF EXISTS SYNCS CASCADE;
	CREATE SCHEMA IF NOT EXISTS linq;

  CREATE TABLE IF NOT EXISTS linq.SYNCS(
    DATE_MAKE TIMESTAMP DEFAULT NOW(),
    DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    _IDT VARCHAR(80) DEFAULT '-1',
    ACTION VARCHAR(80) DEFAULT '',
    _ID VARCHAR(80) DEFAULT '-1',
    _SYNC BOOLEAN DEFAULT FALSE,    
    INDEX SERIAL,
    PRIMARY KEY (TABLE_SCHEMA, TABLE_NAME, _IDT)
  );  
  CREATE INDEX IF NOT EXISTS SYNCS_INDEX_IDX ON linq.SYNCS(INDEX);
  CREATE INDEX IF NOT EXISTS SYNCS_TABLE_SCHEMA_IDX ON linq.SYNCS(TABLE_SCHEMA);
  CREATE INDEX IF NOT EXISTS SYNCS_TABLE_NAME_IDX ON linq.SYNCS(TABLE_NAME);
  CREATE INDEX IF NOT EXISTS SYNCS__IDT_IDX ON linq.SYNCS(_IDT);
  CREATE INDEX IF NOT EXISTS SYNCS_ACTION_IDX ON linq.SYNCS(ACTION);
  CREATE INDEX IF NOT EXISTS SYNCS__ID_IDX ON linq.SYNCS(_ID);
  CREATE INDEX IF NOT EXISTS SYNCS__SYNC_IDX ON linq.SYNCS(_SYNC);

  CREATE OR REPLACE FUNCTION linq.SYNC_INSERT()
  RETURNS
    TRIGGER AS $$
  DECLARE
   CHANNEL VARCHAR(250);
  BEGIN
    IF NEW._IDT = '-1' THEN
      NEW._IDT = uuid_generate_v4();

      INSERT INTO linq.SYNCS(TABLE_SCHEMA, TABLE_NAME, _IDT, ACTION, _ID)
      VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, NEW._IDT, TG_OP, uuid_generate_v4());

      PERFORM pg_notify(
      'sync',
      json_build_object(
        'option', TG_OP,        
        '_idt', NEW._IDT
      )::text
      );

      CHANNEL = TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
      PERFORM pg_notify(
      CHANNEL,
      json_build_object(
        'option', TG_OP,
        '_idt', NEW._IDT
      )::text
      );
    END IF;

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION linq.SYNC_UPDATE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    CHANNEL VARCHAR(250);
  BEGIN
    IF NEW._IDT = '-1' AND OLD._IDT != '-1' THEN
      NEW._IDT = OLD._IDT;
    ELSE
     IF NEW._IDT = '-1' THEN
       NEW._IDT = uuid_generate_v4();
     END IF;
     INSERT INTO linq.SYNCS(TABLE_SCHEMA, TABLE_NAME, _IDT, ACTION, _ID)
     VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, NEW._IDT, TG_OP, uuid_generate_v4())
		 ON CONFLICT(TABLE_SCHEMA, TABLE_NAME, _IDT) DO UPDATE SET
     DATE_UPDATE = NOW(),
     ACTION = TG_OP,
     _SYNC = FALSE,
     _ID = uuid_generate_v4();

     PERFORM pg_notify(
     'sync',
     json_build_object(
       'option', TG_OP,
       '_idt', NEW._IDT
     )::text
     );

     CHANNEL = TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
     PERFORM pg_notify(
     CHANNEL,
     json_build_object(
       'option', TG_OP,
       '_idt', NEW._IDT  
     )::text
     );
    END IF; 

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION linq.SYNC_DELETE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    VINDEX INTEGER;
    CHANNEL VARCHAR(250);
  BEGIN
    SELECT INDEX INTO VINDEX
    FROM linq.SYNCS
    WHERE TABLE_SCHEMA = TG_TABLE_SCHEMA
    AND TABLE_NAME = TG_TABLE_NAME
    AND _IDT = OLD._IDT
    LIMIT 1;
    IF FOUND THEN
      UPDATE linq.SYNCS SET
      DATE_UPDATE = NOW(),
      ACTION = TG_OP,
      _SYNC = FALSE,
      _ID = uuid_generate_v4()
      WHERE INDEX = VINDEX;
      
      PERFORM pg_notify(
      'sync',
      json_build_object(
        'option', TG_OP,
        '_idt', OLD._IDT
      )::text
      );

      CHANNEL = TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
      PERFORM pg_notify(
      CHANNEL,
      json_build_object(
        'option', TG_OP,
        '_idt', OLD._IDT
      )::text
      );      
    END IF;

  RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;
	`
}

// ddlRecicling return sql recicling ddl
func ddlRecicling() string {
	return `
	-- DROP TABLE IF EXISTS linq.RECYCLING CASCADE;
	CREATE SCHEMA IF NOT EXISTS linq;

  CREATE TABLE IF NOT EXISTS linq.RECYCLING(
		TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    _IDT VARCHAR(80) DEFAULT '-1',
    INDEX SERIAL,
		PRIMARY KEY(TABLE_SCHEMA, TABLE_NAME, _IDT)
	);
  CREATE INDEX IF NOT EXISTS RECYCLING_INDEX_IDX ON linq.RECYCLING(INDEX);
	CREATE INDEX IF NOT EXISTS RECYCLING_TABLE_SCHEMA_IDX ON linq.RECYCLING(TABLE_SCHEMA);
	CREATE INDEX IF NOT EXISTS RECYCLING_TABLE_NAME_IDX ON linq.RECYCLING(TABLE_NAME);
	CREATE INDEX IF NOT EXISTS RECYCLING__IDT_IDX ON linq.RECYCLING(INDEX);

	CREATE OR REPLACE FUNCTION linq.RECYCLING()
  RETURNS
    TRIGGER AS $$
  BEGIN
		IF NEW._STATE != OLD._STATE AND NEW._STATE = '-2' THEN
    	INSERT INTO linq.RECYCLING(TABLE_SCHEMA, TABLE_NAME, _IDT)
    	VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, NEW._IDT);

      PERFORM pg_notify(
      'recycling',
      json_build_object(
        '_idt', NEW._IDT
      )::text
      );
		ELSEIF NEW._STATE != OLD._STATE AND OLD._STATE = '-2' THEN
			DELETE FROM linq.RECYCLING WHERE _IDT=NEW._IDT;
    END IF;

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION ERASE()
  RETURNS
    TRIGGER AS $$
  BEGIN
		DELETE FROM linq.RECYCLING WHERE _IDT=OLD._IDT;
  	RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;
	`
}

func ddlStrucs() string {
	return `
	-- DROP TABLE IF EXISTS linq.STRUCTS CASCADE;
	CREATE SCHEMA IF NOT EXISTS linq;

  CREATE TABLE IF NOT EXISTS linq.STRUCTS(
		TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    SQL TEXT DEFAULT '',
    INDEX SERIAL,
		PRIMARY KEY(TABLE_SCHEMA, TABLE_NAME)
	);
  CREATE INDEX IF NOT EXISTS STRUCTS_INDEX_IDX ON linq.STRUCTS(INDEX);
	CREATE INDEX IF NOT EXISTS STRUCTS_SCHEMA_IDX ON linq.STRUCTS(TABLE_SCHEMA);
	CREATE INDEX IF NOT EXISTS STRUCTS_NAME_IDX ON linq.STRUCTS(TABLE_NAME);

	CREATE OR REPLACE FUNCTION linq.setstruct(
	VSCHEMA VARCHAR(80),
	VNAME VARCHAR(80),
	VSQL TEXT)
	RETURNS INT AS $$
	DECLARE
	 result INT;
	BEGIN
	 INSERT INTO linq.STRUCTS AS A (TABLE_SCHEMA, TABLE_NAME, SQL)
	 SELECT VSCHEMA, VNAME, VSQL
	 ON CONFLICT (TABLE_SCHEMA, TABLE_NAME) DO UPDATE SET
	 SQL = VSQL
	 RETURNING INDEX INTO result;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	`
}

func ddlSeries() string {
	return `
	-- DROP TABLE IF EXISTS linq.SERIES CASCADE;
	CREATE SCHEMA IF NOT EXISTS linq;

  CREATE TABLE IF NOT EXISTS linq.SERIES(		
		SERIE VARCHAR(250) DEFAULT '',
		VALUE BIGINT DEFAULT 0,
		PRIMARY KEY(SERIE)
	);
	
	CREATE OR REPLACE FUNCTION linq.nextserie(tag VARCHAR(250))
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 INSERT INTO linq.SERIES AS A (SERIE, VALUE)
	 SELECT tag, 1
	 ON CONFLICT (SERIE) DO UPDATE SET
	 VALUE = A.VALUE + 1
	 RETURNING VALUE INTO result;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION linq.setserie(tag VARCHAR(250), val BIGINT)
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 INSERT INTO linq.SERIES AS A (SERIE, VALUE)
	 SELECT tag, val
	 ON CONFLICT (SERIE) DO UPDATE SET
	 VALUE = val
	 WHERE A.VALUE < val
	 RETURNING VALUE INTO result;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION linq.currserie(tag VARCHAR(250))
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 SELECT VALUE INTO result
	 FROM linq.SERIES
	 WHERE SERIE = tag LIMIT 1;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	`
}

// ddlFuntions return sql funcitions ddl to support a models
func ddlFuntions() string {
	return `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE SCHEMA IF NOT EXISTS linq;
	
	CREATE OR REPLACE FUNCTION linq.create_constraint_if_not_exists(
	s_name text,
	t_name text,
	c_name text,
	constraint_sql text) 
	RETURNS void AS $$
	BEGIN
		IF NOT EXISTS(
		SELECT constraint_name 
		FROM information_schema.table_constraints 
		WHERE UPPER(table_schema)=UPPER(s_name)
		AND UPPER(table_name)=UPPER(t_name)
		AND UPPER(constraint_name)=UPPER(c_name)) THEN
		 execute constraint_sql;
		END IF;
	END;
	$$ LANGUAGE 'plpgsql';
	`
}

// ddlDedault return sql default values
func ddlDefault(col *linq.Column) string {
	var result string
	switch col.TypeData {
	case linq.TpKey:
		result = `'-1'`
	case linq.TpText:
		result = `''`
	case linq.TpMemo:
		result = `''`
	case linq.TpNumber:
		result = `0`
	case linq.TpDate:
		result = `NOW()`
	case linq.TpCheckbox:
		result = `FALSE`
	case linq.TpRelation:
		result = `''`
	case linq.TpRollup:
		result = `''`
	case linq.TpCreatedTime:
		result = `NOW()`
	case linq.TpCreatedBy:
		result = `'{ "_id": "", "name": "" }'`
	case linq.TpLastEditedTime:
		result = `NOW()`
	case linq.TpLastEditedBy:
		result = `'{ "_id": "", "name": "" }'`
	case linq.TpStatus:
		result = `'{ "_id": "0", "main": "State", "name": "Activo" }'`
	case linq.TpPerson:
		result = `'{ "_id": "", "name": "" }'`
	case linq.TpFile:
		result = `''`
	case linq.TpURL:
		result = `''`
	case linq.TpEmail:
		result = `''`
	case linq.TpPhone:
		result = `''`
	case linq.TpFormula:
		result = `''`
	case linq.TpSelect:
		result = `''`
	case linq.TpMultiSelect:
		result = `''`
	case linq.TpJson:
		result = `'{}'`
	case linq.TpArray:
		result = `'[]'`
	case linq.TpSerie:
		result = `0`
	default:
		val := col.Default
		result = strs.Format(`%v`, et.Quote(val))
	}

	return strs.Append("DEFAULT", result, " ")
}

// ddlType return sql type ddl
func ddlType(col *linq.Column) string {
	switch col.TypeData {
	case linq.TpKey, linq.TpRelation, linq.TpRollup, linq.TpStatus, linq.TpPhone, linq.TpSelect, linq.TpMultiSelect:
		return "VARCHAR(80)"
	case linq.TpMemo:
		return "TEXT"
	case linq.TpNumber:
		return "DECIMAL(18, 2)"
	case linq.TpDate:
		return "TIMESTAMP"
	case linq.TpCheckbox:
		return "BOOLEAN"
	case linq.TpCreatedTime:
		return "TIMESTAMP"
	case linq.TpCreatedBy:
		return "JSONB"
	case linq.TpLastEditedTime:
		return "TIMESTAMP"
	case linq.TpLastEditedBy:
		return "JSONB"
	case linq.TpPerson:
		return "JSONB"
	case linq.TpFile:
		return "JSONB"
	case linq.TpURL:
		return "TEXT"
	case linq.TpFormula:
		return "JSONB"
	case linq.TpJson:
		return "JSONB"
	case linq.TpArray:
		return "JSONB"
	case linq.TpSerie:
		return "BIGINT"
	default:
		return "VARCHAR(250)"
	}
}

// ddlschema return sql schema ddl
func ddlSchema(schema *linq.Schema) string {
	return strs.Format(`CREATE SCHEMA IF NOT EXISTS %s;`, schema.Name)
}

// ddlColumn return sql column ddl
func ddlColumn(col *linq.Column) string {
	var result string
	var def string

	result = ddlDefault(col)
	def = ddlType(col)
	result = strs.Append(def, result, " ")
	result = strs.Append(col.Up(), result, " ")

	return result
}

// ddlIndex return sql index ddl
func ddlIndex(col *linq.Column) string {
	name := strs.Format(`%v_%v_IDX`, strs.Uppcase(col.Table()), col.Up())
	name = strs.Replace(name, "-", "_")
	name = strs.Replace(name, ".", "_")
	return strs.Format(`CREATE INDEX IF NOT EXISTS %v ON %v(%v);`, name, strs.Uppcase(col.Table()), col.Up())
}

// ddlUnique return sql unique index ddl
func ddlUnique(col *linq.Column) string {
	name := strs.Format(`%v_%v_IDX`, strs.Uppcase(col.Table()), col.Up())
	name = strs.Replace(name, "-", "_")
	name = strs.Replace(name, ".", "_")
	return strs.Format(`CREATE UNIQUE INDEX IF NOT EXISTS %v ON %v(%v);`, name, strs.Uppcase(col.Table()), col.Up())
}

// ddlPrimaryKey return sql primary key ddl
func ddlPrimaryKey(col *linq.Column) string {
	schema := strs.Uppcase(col.Model.Schema.Name)
	key := strs.Replace(col.Table(), ".", "_")
	key = strs.Replace(key, "-", "_") + "_pkey"
	key = strs.Lowcase(key)
	def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s PRIMARY KEY (%s);`, strs.Uppcase(col.Table()), key, strings.Join(col.PrimaryKeys(), ", "))
	return strs.Format(`SELECT linq.create_constraint_if_not_exists('%s', '%s', '%s', '%s');`, schema, col.Model.Name, key, def)
}

// ddlForeignKeys return ForeignKey ddl
func ddlForeignKeys(model *linq.Model) string {
	var result string
	for _, ref := range model.ForeignKey {
		schema := strs.Uppcase(model.Schema.Name)
		key := strs.Replace(model.Table, ".", "_")
		key = strs.Replace(key, "-", "_") + "_pkey"
		key = strs.Lowcase(key)
		def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s);`, strs.Uppcase(model.Table), key, strings.Join(ref.ForeignKey, ", "), ref.Parent.Table, strings.Join(ref.ParentKey, ", "))
		def = strs.Format(`SELECT linq.create_constraint_if_not_exists('%s', '%s', '%s', '%s');`, schema, model.Name, key, def)
		result = strs.Append(result, def, "\n")
	}

	return result
}

// ddlTable return table ddl
func ddlTable(model *linq.Model) string {
	var result string
	var columns string
	var indexs string
	for _, col := range model.Columns {
		if col.TypeColumn == linq.TpColumn {
			def := ddlColumn(col)
			columns = strs.Append(def, columns, ",\n")
			if col.PrimaryKey {
				def = ddlPrimaryKey(col)
				indexs = strs.Append(def, indexs, "\n")
			} else if col.Unique {
				def = ddlUnique(col)
				indexs = strs.Append(def, indexs, "\n")
			} else if col.Indexed {
				def = ddlIndex(col)
				indexs = strs.Append(def, indexs, "\n")
			}
		}
	}
	schema := ddlSchema(model.Schema)
	result = strs.Append(result, schema, "\n")
	table := strs.Format("CREATE TABLE IF NOT EXISTS %s (\n%s);", model.Table, columns)
	result = strs.Append(result, table, "\n")
	result = strs.Append(result, indexs, "\n")
	foreign := ddlForeignKeys(model)
	result = strs.Append(result, foreign, "\n")

	return result
}
