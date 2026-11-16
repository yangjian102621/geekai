package library

import "context"

func (r *remote) Collections(ctx context.Context) ([]Collection, error) {
	out := envelope[[]Collection]{}
	err := r.get(ctx, "/api/v1/knowledge-bases", &out)
	return out.Data, out.check(err)
}

func (r *remote) Collection(ctx context.Context, id string) (Collection, error) {
	if !ValidID(id) {
		return Collection{}, ErrNotFound
	}
	out := envelope[Collection]{}
	err := r.get(ctx, "/api/v1/knowledge-bases/"+id, &out)
	if err = out.check(err); err != nil {
		return Collection{}, err
	}
	if out.Data.ID != id {
		return Collection{}, ErrNotFound
	}
	return out.Data, nil
}
