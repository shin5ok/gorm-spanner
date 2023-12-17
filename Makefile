SPANNER_INSTANCE := test-instance
SPANNER_DATABASE := game
SPANNER_STRING := projects/$(GOOGLE_CLOUD_PROJECT)/instances/$(SPANNER_INSTANCE)/databases/$(SPANNER_DATABASE)
SPANNER_CONFIG := regional-us-east5


.PHONY: create-spanner-instance
create-spanner-instance:
	gcloud spanner instances create --config=$(SPANNER_CONFIG) \
    --instance-type=free-instance --description="test-free-instance" $(SPANNER_INSTANCE)

.PHONY: create-spanner-database
create-spanner-database:
	gcloud spanner databases create $(SPANNER_DATABASE) --instance=$(SPANNER_INSTANCE)

.PHONY: create-schema
create-schema:
	@echo "Creating schemas to Cloud Spanner databse $(SPANNER_DATABASE) at $(SPANNER_DATABASE)"
	for schema in schemas/*ddl.sql schemas/*dml.sql ; do spanner-cli -i $(SPANNER_INSTANCE) -d $(SPANNER_DATABASE) -p $(GOOGLE_CLOUD_PROJECT) < $${schema} ; done
