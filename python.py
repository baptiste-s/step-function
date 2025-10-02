import json
import psycopg2
import os
import time
import random
from datetime import datetime

def lambda_handler(event, context):
    """
    Lambda worker qui effectue des inserts unitaires sur PostgreSQL
    
    Paramètres attendus dans event:
    - num_inserts: nombre d'inserts à effectuer
    - table_name: nom de la table (défaut: test_data)
    - worker_id: identifiant du worker
    """
    
    # Récupération des paramètres
    num_inserts = event.get('num_inserts', 1000)
    table_name = event.get('table_name', 'test_data')
    worker_id = event.get('worker_id', context.request_id)
    
    # Connexion à la base de données
    db_host = os.environ['DB_HOST']
    db_name = os.environ['DB_NAME']
    db_user = os.environ['DB_USER']
    db_password = os.environ['DB_PASSWORD']
    db_port = os.environ.get('DB_PORT', '5432')
    
    start_time = time.time()
    successful_inserts = 0
    failed_inserts = 0
    
    try:
        # Connexion à PostgreSQL
        conn = psycopg2.connect(
            host=db_host,
            database=db_name,
            user=db_user,
            password=db_password,
            port=db_port,
            connect_timeout=10
        )
        
        cursor = conn.cursor()
        
        # Création de la table si elle n'existe pas
        create_table_query = f"""
        CREATE TABLE IF NOT EXISTS {table_name} (
            id SERIAL PRIMARY KEY,
            worker_id VARCHAR(100),
            timestamp TIMESTAMP,
            value INTEGER,
            text_data VARCHAR(255),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
        """
        cursor.execute(create_table_query)
        conn.commit()
        
        # Effectuer les inserts unitaires
        for i in range(num_inserts):
            try:
                insert_query = f"""
                INSERT INTO {table_name} (worker_id, timestamp, value, text_data)
                VALUES (%s, %s, %s, %s)
                """
                
                cursor.execute(insert_query, (
                    worker_id,
                    datetime.now(),
                    random.randint(1, 1000),
                    f"Test data {random.randint(1, 10000)}"
                ))
                conn.commit()
                successful_inserts += 1
                
            except Exception as e:
                failed_inserts += 1
                print(f"Erreur lors de l'insert {i}: {str(e)}")
                conn.rollback()
        
        cursor.close()
        conn.close()
        
        end_time = time.time()
        duration = end_time - start_time
        
        return {
            'statusCode': 200,
            'body': json.dumps({
                'worker_id': worker_id,
                'successful_inserts': successful_inserts,
                'failed_inserts': failed_inserts,
                'duration_seconds': round(duration, 2),
                'inserts_per_second': round(successful_inserts / duration, 2) if duration > 0 else 0,
                'start_time': start_time,
                'end_time': end_time
            })
        }
        
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps({
                'worker_id': worker_id,
                'error': str(e),
                'successful_inserts': successful_inserts,
                'failed_inserts': failed_inserts
            })
        }
