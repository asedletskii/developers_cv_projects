DELETE FROM resume_contact_hidden AS hidden
USING resume_contacts AS contacts
WHERE hidden.contacts_id = contacts.id
  AND (
      (hidden.value = 'phone' AND contacts.phone IS NULL)
      OR (hidden.value = 'telegram' AND contacts.telegram IS NULL)
      OR (hidden.value = 'email' AND contacts.email IS NULL)
  );
