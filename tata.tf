def insert_parallel_mode(db_host, db_name, db_user, db_password, db_port,
                        table_name, worker_id, num_inserts, num_threads):
    """Mode parallèle - plusieurs threads avec connexions séparées (TRÈS RAPIDE)"""
    
    def insert_chunk(thread_id, chunk_size):
        """Fonction exécutée par chaque thread"""
        # CHAQUE THREAD A SA PROPRE CONNEXION
        conn = psycopg2.connect(
            host=db_host, database=db_name, user=db_user,
            password=db_password, port=db_port
        )
        cursor = conn.cursor()
        
        # ... fait ses inserts ...
        
        conn.close()
        return successful, failed
    
    # Lance 10 threads en parallèle (par défaut)
    with ThreadPoolExecutor(max_workers=num_threads) as executor:
        futures = []
        for i in range(num_threads):
            futures.append(executor.submit(insert_chunk, i, size))
        
        # Attend que tous les threads terminent
        for future in as_completed(futures):
            successful, failed = future.result()
