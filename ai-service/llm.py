import os
import json
from gigachat import GigaChat
from gigachat.models import Chat, Messages, MessagesRole

API_KEY = os.environ.get("GIGACHAT_API_KEY")


def find_entities_llm(text: str) -> list[dict]:
    if not API_KEY:
        return []

    prompt = (
        'Найди персональные данные в тексте. Верни только JSON массив, без пояснений.\n\n'
        'Формат: [{"type": "тип", "value": "значение"}]\n'
        'Типы: phone, email, inn, snils, passport, address, name\n\n'
        'Текст:\n' + text
    )

    try:
        with GigaChat(credentials=API_KEY, verify_ssl_certs=False) as giga:
            response = giga.chat(Chat(
                messages=[Messages(role=MessagesRole.USER, content=prompt)]
            ))
        return json.loads(response.choices[0].message.content)
    except Exception:
        return []