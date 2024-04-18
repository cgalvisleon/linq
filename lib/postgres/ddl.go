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
	CREATE OR REPLACE FUNCTION SYNC_INSERT()
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

  CREATE OR REPLACE FUNCTION SYNC_UPDATE()
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

  CREATE OR REPLACE FUNCTION SYNC_DELETE()
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

  CREATE TABLE IF NOT EXISTS SYNCS(
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
  CREATE INDEX IF NOT EXISTS SYNCS_INDEX_IDX ON SYNCS(INDEX);
  CREATE INDEX IF NOT EXISTS SYNCS_TABLE_SCHEMA_IDX ON SYNCS(TABLE_SCHEMA);
  CREATE INDEX IF NOT EXISTS SYNCS_TABLE_NAME_IDX ON SYNCS(TABLE_NAME);
  CREATE INDEX IF NOT EXISTS SYNCS__IDT_IDX ON SYNCS(_IDT);
  CREATE INDEX IF NOT EXISTS SYNCS_ACTION_IDX ON SYNCS(ACTION);
  CREATE INDEX IF NOT EXISTS SYNCS__ID_IDX ON SYNCS(_ID);
  CREATE INDEX IF NOT EXISTS SYNCS__SYNC_IDX ON SYNCS(_SYNC);

  CREATE OR REPLACE FUNCTION SYNC_INSERT()
  RETURNS
    TRIGGER AS $$
  DECLARE
   CHANNEL VARCHAR(250);
  BEGIN
    IF NEW._IDT = '-1' THEN
      NEW._IDT = uuid_generate_v4();

      INSERT INTO SYNCS(TABLE_SCHEMA, TABLE_NAME, _IDT, ACTION, _ID)
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

  CREATE OR REPLACE FUNCTION SYNC_UPDATE()
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
     INSERT INTO SYNCS(TABLE_SCHEMA, TABLE_NAME, _IDT, ACTION, _ID)
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

  CREATE OR REPLACE FUNCTION SYNC_DELETE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    VINDEX INTEGER;
    CHANNEL VARCHAR(250);
  BEGIN
    SELECT INDEX INTO VINDEX
    FROM SYNCS
    WHERE TABLE_SCHEMA = TG_TABLE_SCHEMA
    AND TABLE_NAME = TG_TABLE_NAME
    AND _IDT = OLD._IDT
    LIMIT 1;
    IF FOUND THEN
      UPDATE SYNCS SET
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
	-- DROP TABLE IF EXISTS RECYCLING CASCADE;

  CREATE TABLE IF NOT EXISTS RECYCLING(
		TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    _IDT VARCHAR(80) DEFAULT '-1',
    INDEX SERIAL,
		PRIMARY KEY(TABLE_SCHEMA, TABLE_NAME, _IDT)
	);
  CREATE INDEX IF NOT EXISTS RECYCLING_INDEX_IDX ON RECYCLING(INDEX);
	CREATE INDEX IF NOT EXISTS RECYCLING_TABLE_SCHEMA_IDX ON RECYCLING(TABLE_SCHEMA);
	CREATE INDEX IF NOT EXISTS RECYCLING_TABLE_NAME_IDX ON RECYCLING(TABLE_NAME);
	CREATE INDEX IF NOT EXISTS RECYCLING__IDT_IDX ON RECYCLING(INDEX);

	CREATE OR REPLACE FUNCTION RECYCLING()
  RETURNS
    TRIGGER AS $$
  BEGIN
		IF NEW._STATE != OLD._STATE AND NEW._STATE = '-2' THEN
    	INSERT INTO RECYCLING(TABLE_SCHEMA, TABLE_NAME, _IDT)
    	VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, NEW._IDT);

      PERFORM pg_notify(
      'recycling',
      json_build_object(
        '_idt', NEW._IDT
      )::text
      );
		ELSEIF NEW._STATE != OLD._STATE AND OLD._STATE = '-2' THEN
			DELETE FROM RECYCLING WHERE _IDT=NEW._IDT;
    END IF;

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION ERASE()
  RETURNS
    TRIGGER AS $$
  BEGIN
		DELETE FROM RECYCLING WHERE _IDT=OLD._IDT;
  	RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;
	`
}

func ddlStrucs() string {
	return `
	-- DROP TABLE IF EXISTS STRUCTS CASCADE;

  CREATE TABLE IF NOT EXISTS STRUCTS(
		TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    SQL TEXT DEFAULT '',
    INDEX SERIAL,
		PRIMARY KEY(TABLE_SCHEMA, TABLE_NAME)
	);
  CREATE INDEX IF NOT EXISTS STRUCTS_INDEX_IDX ON STRUCTS(INDEX);
	CREATE INDEX IF NOT EXISTS STRUCTS_SCHEMA_IDX ON STRUCTS(TABLE_SCHEMA);
	CREATE INDEX IF NOT EXISTS STRUCTS_NAME_IDX ON STRUCTS(TABLE_NAME);

	CREATE OR REPLACE FUNCTION setstruct(
	VSCHEMA VARCHAR(80),
	VNAME VARCHAR(80),
	VSQL TEXT)
	RETURNS INT AS $$
	DECLARE
	 result INT;
	BEGIN
	 INSERT INTO STRUCTS AS A (TABLE_SCHEMA, TABLE_NAME, SQL)
	 SELECT VSCHEMA, VNAME, VSQL
	 ON CONFLICT (TABLE_SCHEMA, TABLE_NAME) DO UPDATE SET
	 SQL = VSQL
	 RETURNING INDEX INTO result;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	`
}

// ddlFuntions return sql funcitions ddl to support a models
func ddlFuntions() string {
	return `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	
	CREATE OR REPLACE FUNCTION create_constraint_if_not_exists(
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
	switch col.Default {
	case linq.DefUuid:
		result = `'-1'`
	case linq.DefInt:
		result = `0`
	case linq.DefInt64:
		result = `0`
	case linq.DefFloat:
		result = `0.0`
	case linq.DefBool:
		result = `FALSE`
	case linq.DefNow:
		result = `NOW()`
	case linq.DefJson:
		result = `'{}'`
	case linq.DefArray:
		result = `'[]'`
	case linq.DefObject:
		result = `'{}'`
	case linq.DefSerie:
		result = `0`
	default:
		val := col.Default.Value()
		result = strs.Format(`%v`, et.Unquote(val))
	}

	return strs.Append("DEFAULT", result, " ")
}

// ddlType return sql type ddl
func ddlType(col *linq.Column) string {
	switch col.TypeData {
	case linq.TpUUId:
		return "VARCHAR(80)"
	case linq.TpInt:
		return "INT"
	case linq.TpInt64:
		return "BIGINT"
	case linq.TpFloat:
		return "DECIMAL(18,2)"
	case linq.TpBool:
		return "BOOLEAN"
	case linq.TpDateTime:
		return "TIMESTAMP"
	case linq.TpTimeStamp:
		return "TIMESTAMP"
	case linq.TpJson:
		return "JSONB"
	case linq.TpArray:
		return "JSONB"
	case linq.TpSerie:
		return "BIGINT"
	case linq.TpText:
		return "TEXT"
	default:
		return "VARCHAR(255)"
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
	pkey := strs.Replace(col.Table(), ".", "_")
	pkey = strs.Replace(pkey, "-", "_") + "_pkey"
	pkey = strs.Lowcase(pkey)
	def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s PRIMARY KEY (%s);`, strs.Uppcase(col.Table()), pkey, strings.Join(col.PrimaryKeys(), ", "))
	return strs.Format(`SELECT create_constraint_if_not_exists('%s', '%s', '%s', '%s');`, col.Schema.Name, col.Model.Name, pkey, def)
}

// ddlForeignKeys return ForeignKey ddl
func ddlForeignKeys(model *linq.Model) string {
	var result string
	for _, ref := range model.ForeignKey {
		def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s);`, strs.Uppcase(model.Table), ref.Name, strings.Join(ref.ForeignKey, ", "), ref.ParentModel.Table, strings.Join(ref.ParentKey, ", "))
		def = strs.Format(`SELECT create_constraint_if_not_exists('%s', '%s', '%s', '%s');`, model.Schema.Name, model.Name, ref.Name, def)
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
