import os, json
from html.parser import HTMLParser

class PageParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.title = ''
        self.in_title = False
        self.sidebar = False
        self.right = False
        self.in_a = False
        self.a_href = None
        self.a_text = ''
        self.nav = []
        self.right_links = []
        self.current_acc = None
        self.current_section = None
        self.in_h3 = False
        self.acc = []
    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == 'title':
            self.in_title = True
        if tag == 'aside' and attrs.get('class') == 'sidebar':
            self.sidebar = True
        if tag == 'aside' and attrs.get('class') == 'right-sidebar':
            self.right = True
        if tag == 'a':
            self.in_a = True
            self.a_href = attrs.get('href')
            self.a_text = ''
        if tag == 'div':
            cls = attrs.get('class', '')
            if cls.startswith('accordion-item'):
                self.current_acc = {'header': None, 'items': []}
            if cls.startswith('ticket-overview') and self.current_acc is not None:
                self.current_section = {'text': '', 'input_id': None, 'label_class': None}
        if tag == 'input' and self.current_section is not None:
            self.current_section['input_id'] = attrs.get('id')
        if tag == 'label' and self.current_section is not None:
            self.current_section['label_class'] = attrs.get('class')
        if tag == 'h3':
            self.in_h3 = True
    def handle_endtag(self, tag):
        if tag == 'title':
            self.in_title = False
        if tag == 'aside' and self.sidebar:
            self.sidebar = False
        if tag == 'aside' and self.right:
            self.right = False
        if tag == 'a' and self.in_a:
            item = {'href': self.a_href, 'text': self.a_text.strip()}
            if self.sidebar and not self.right:
                self.nav.append(item)
            elif self.right:
                self.right_links.append(item)
            self.in_a = False
        if tag == 'div' and self.current_section is not None:
            if self.current_acc is not None:
                self.current_acc['items'].append(self.current_section)
            self.current_section = None
        if tag == 'h3':
            self.in_h3 = False
    def handle_data(self, data):
        if self.in_title:
            self.title += data.strip()
        if self.in_a:
            self.a_text += data
        if self.in_h3 and self.current_section is not None:
            self.current_section['text'] += data.strip()
        if self.in_h3 and self.current_section is None and self.current_acc is not None:
            text = data.strip()
            if text:
                self.current_acc['header'] = text
                self.acc.append(self.current_acc)
                self.current_acc = None

root = 'html'
all_pages = []
for dirpath, dirnames, filenames in os.walk(root):
    for fname in sorted(filenames):
        if fname.endswith('.html'):
            path = os.path.join(dirpath, fname)
            with open(path, 'r', encoding='utf-8') as f:
                html = f.read()
            parser = PageParser()
            parser.feed(html)
            all_pages.append({'path': path, 'title': parser.title, 'nav': parser.nav, 'right': parser.right_links, 'acc': parser.acc})
print(json.dumps(all_pages, indent=2, ensure_ascii=False))
