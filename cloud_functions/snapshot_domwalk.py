from google.cloud import bigquery

def snapshot_dataset(data, context):
    """
    Cloud Function to snapshot all tables in a BigQuery dataset.

    Args:
        data (dict): The event payload.
        context (google.cloud.functions.Context): Metadata for the event.
    """

    client = bigquery.Client()

    # Replace with your source dataset ID
    source_dataset_id = "your-source-dataset"

    # Replace with your snapshot dataset ID
    snapshot_dataset_id = "your-snapshot-dataset"

    source_dataset_ref = client.dataset(source_dataset_id)
    snapshot_dataset_ref = client.dataset(snapshot_dataset_id)

    for table in client.list_tables(source_dataset_ref):
        source_table_ref = source_dataset_ref.table(table.table_id)
        snapshot_table_id = f"{table.table_id}_{datetime.datetime.now().strftime('%Y%m%d%H%M%S')}"
        snapshot_table_ref = snapshot_dataset_ref.table(snapshot_table_id)

        job_config = bigquery.CopyJobConfig()
        job_config.write_disposition = "WRITE_TRUNCATE"  # Overwrite existing snapshots

        job = client.copy_table(source_table_ref, snapshot_table_ref, job_config=job_config)
        job.result()  # Wait for the job to complete

        print(f"Snapshot created for table {source_table_ref} to {snapshot_table_ref}")