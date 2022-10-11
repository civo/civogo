package civogo

import (
	"testing"
)

func TestListActions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/actions": `{"page":1,"per_page":5,"pages":1,"items":[{"id":267531707,"created_at":"2022-10-10T16:30:11Z","updated_at":"2022-10-10T16:30:11Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"","type":"cluster-create","details":"Created a new cluster called cluster-kubectl","related_id":"259e96c5-ecd4-43c6-be35-d21dcf05b650","related_type":"cluster","debug":false},{"id":267531696,"created_at":"2022-10-10T16:29:50Z","updated_at":"2022-10-10T16:29:50Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"6b6e9801-8bcd-4049-84d7-028c8e748f58","type":"cluster-delete","details":"Deleted cluster : cluster-kubectl-c387-e36f8c","related_id":"7b107c00-28b1-49b8-b609-c850aeb2d72e","related_type":"cluster","debug":false},{"id":267531012,"created_at":"2022-10-10T12:08:09Z","updated_at":"2022-10-10T12:08:09Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"","type":"cluster-create","details":"Created a new cluster called cluster-kubectl","related_id":"7b107c00-28b1-49b8-b609-c850aeb2d72e","related_type":"cluster","debug":false},{"id":267527196,"created_at":"2022-10-09T12:40:37Z","updated_at":"2022-10-09T12:40:37Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"6b6e9801-8bcd-4049-84d7-028c8e748f58","type":"network-delete","details":"Deleted a network called cust-test-kubectl-eaef1dd6-5b25-58204a62","related_id":"0d7be2cc-515d-44c2-8b84-1dbd686806eb","related_type":"network","debug":false},{"id":267527195,"created_at":"2022-10-09T12:40:22Z","updated_at":"2022-10-09T12:40:22Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"6b6e9801-8bcd-4049-84d7-028c8e748f58","type":"volume-delete","details":"Deleted volume : test-server-665b-21da40","related_id":"0b5313a7-9b59-43c0-b749-0a0195407a62","related_type":"volume","debug":false}]}`,
	})
	defer server.Close()

	actionListRequest := &ActionListRequest{}
	allActions, err := client.ListActions(actionListRequest)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if len(allActions.Items) != 5 {
		t.Errorf("Expected %d, got %d", 5, len(allActions.Items))
	}
}

func TestListActionsWithFilter(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/actions": `{"page":1,"per_page":5,"pages":1,"items":[{"id":267531707,"created_at":"2022-10-10T16:30:11Z","updated_at":"2022-10-10T16:30:11Z","account_id":"eaef1dd6-1cec-4d9c-8480-96452bd94dea","user_id":"","type":"cluster-create","details":"Created a new cluster called cluster-kubectl","related_id":"259e96c5-ecd4-43c6-be35-d21dcf05b650","related_type":"cluster","debug":false}]}`,
	})
	defer server.Close()

	actionListRequest := &ActionListRequest{
		RelatedID: "259e96c5-ecd4-43c6-be35-d21dcf05b650",
	}
	allActions, err := client.ListActions(actionListRequest)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if len(allActions.Items) != 1 {
		t.Errorf("Expected %d, got %d", 1, len(allActions.Items))
	}
}
