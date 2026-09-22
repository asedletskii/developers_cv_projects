package mock

import (
	"context"
	"time"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

const (
	PersonOneID = "550e8400-e29b-41d4-a716-446655440001"
	PersonTwoID = "550e8400-e29b-41d4-a716-446655440002"
)

type PersonReader struct {
	people map[uuid.UUID]queries.Person
}

func NewPersonReader() *PersonReader {
	first := newFirstPerson()
	second := newSecondPerson(first)

	return &PersonReader{
		people: map[uuid.UUID]queries.Person{
			uuid.MustParse(PersonOneID): first,
			uuid.MustParse(PersonTwoID): second,
		},
	}
}

func (r *PersonReader) GetPerson(_ context.Context, query queries.GetPersonQuery) (queries.Person, error) {
	person, ok := r.people[query.ID]
	if !ok {
		return queries.Person{}, queries.ErrPersonNotFound
	}

	person = clonePerson(person)
	person.Greeting = localizedGreeting(query.ID, query.Lang, person.Greeting)

	if person.Availability != nil {
		person.Availability.Note = localizedAvailabilityNote(query.ID, query.Lang)
	}

	person.ContactsText = localizedContactsText(query.ID, query.Lang)

	updatedAt := time.Now().UTC().Add(-72 * time.Hour)
	person.UpdatedAt = &updatedAt

	return person, nil
}

func clonePerson(person queries.Person) queries.Person {
	person.Sections = append([]queries.SectionSettings{}, person.Sections...)
	person.HiddenContacts = append([]string{}, person.HiddenContacts...)
	person.WorkFormat = append([]string{}, person.WorkFormat...)
	person.Languages = append([]queries.PersonLanguage{}, person.Languages...)
	person.Experience = append([]queries.Experience{}, person.Experience...)
	person.Education = append([]queries.Education{}, person.Education...)
	person.Courses = append([]queries.Course{}, person.Courses...)
	person.Activities = append([]queries.Activity{}, person.Activities...)
	person.HardSkills = append([]string{}, person.HardSkills...)
	person.SoftSkills = append([]string{}, person.SoftSkills...)
	person.Technologies = append([]queries.Technology{}, person.Technologies...)
	person.Projects = append([]queries.ProjectSummary{}, person.Projects...)
	if person.Availability != nil {
		availability := *person.Availability
		person.Availability = &availability
	}
	if person.ContactsText != nil {
		contactsText := *person.ContactsText
		person.ContactsText = &contactsText
	}

	return person
}
