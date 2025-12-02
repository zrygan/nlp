"""
Filipino Naturalization CFG Parser
Based on Filipino orthography rules from linguistic research
Applies phonological rules as documented in Filipino language standardization
"""

import re
from typing import List, Tuple, Dict, Optional

class FilipinoCFGParser:
    def __init__(self):
        # Core Filipino alphabet
        self.core_consonants = {'b', 'k', 'd', 'g', 'h', 'l', 'm', 'n', 'ng', 'p', 'r', 's', 't', 'w', 'y'}
        self.vowels = {'a', 'e', 'i', 'o', 'u'}
        
        # Phonological mapping rules (Section 2.1.4 of the document)
        # These are the official Filipino naturalization rules
        self.phoneme_rules = [
            ('ph', 'f'),   # phone -> pon
            ('ps', 's'),   # psychology -> sikolodi
            ('v', 'b'),    # verde -> berde
            ('j', 'dy'),   # jeep -> dyip
            ('z', 's'),    # zipon -> sipon
            ('ñ', 'ny'),   # baño -> banyo
            ('qu', 'kuw'),  # queen -> kwin
        ]
        
        # Special C handling: /s/ sound -> s, /k/ sound -> k
        # gobierno -> gobyerno (c as /k/ before i/e in Spanish -> /g/ in Filipino)
        # centro -> sentro (c as /s/)
        
        # Common Spanish->Filipino transformations
        self.spanish_patterns = [
            (r'é', 'e'),             # teléfono -> teléphono
            (r'o', 'o'),             # zipon -> sipon
            (r'á', 'a'),             # como está -> komo está
            (r'cion$', 'syon'),      # educacion -> edukasyon
            (r'cion$', 'syon'),      # 
            (r'tion$', 'syon'),      # education -> edukasyon
            (r'tion$', 'syon'),      # education -> edukasyon
            (r'gobierno', 'gobyerno'), # gobierno -> gobyerno
            (r'bie', 'biye'),        # gobierno -> gobyerno
            (r'ci([eo])', r'sy\1'),  # gracias -> grasyas
            (r'ce', 'se'),           # centro -> sentro
            (r'ci', 'si'),           # medicina -> medisina
            (r'o e', 'u'),           # como está -> komustá
            (r'^ko', 'ku'),          # como está -> kumustá
        ]
        
        # Vowel interchangeability (Section 2.1.1)
        # e/i and o/u are allophones in Filipino
        self.vowel_shifts = {
            'ee': 'i',
            'oo': 'u',
        }
        
        # Common English->Filipino word patterns
        self.english_patterns = [
            (r'tion$', 'syon'),
            (r'sion$', 'syon'),
            (r'ture$', 'tyur'),
            (r'ter$', 'ter'),
            (r'ble$', 'bol'),
        ]
        
    def apply_phonological_rules(self, word: str, source_lang: str = 'auto') -> str:
        """
        Apply Filipino phonological naturalization rules
        
        Args:
            word: Word to naturalize
            source_lang: 'spanish', 'english', or 'auto'
        """
        result = word.lower()
        
        # Step 1: Apply basic phoneme mappings (works for both Spanish and English)
        for pattern, replacement in self.phoneme_rules:
            result = result.replace(pattern, replacement)
        
        # if y is after x, change to i (e.g., acyclovir -> asayklobir)
        result = re.sub(r'(?<=[x])y', 'i', result)
        
        # Only if starts with x, change to s
        result = re.sub(r'^x', 's', result)
        
        # Every other case, replace with ks
        result = re.sub(r'x', 'ks', result)
        
        # if y is before s, change to i (e.g., acyclovir -> asayklobir)
        result = re.sub(r'y(?=[s])', 'i', result)    
        
        # If y is before l, change to i (e.g., acetylcysteine -> asetilsisteen)
        result = re.sub(r'(?i)y(?=l)', 'i', result)    
        
        # c before e/i -> s (centro -> sentro)
        result = re.sub(r'c([eyi])', r's\1', result)
        
        # c elsewhere -> k (computer -> kompyuter)
        result = result.replace('c', 'k')
        
        # Any ch -> k instead
        result = result.replace('kh', 'k')
        
        # Special handling for 'sy' -> 'sayk' (e.g., cyclosporine -> siklosporine)
        result = re.sub(r'(?i)y', 'ay', result)
        
        # switch y and i if y does not preceed with a
        result = re.sub(r'(?<!a)y', 'i', result)
        
        # Handle anything ending with -ate
        result = re.sub(r'ate$', 'eyt', result)
        
        # Handle anything ending with -ide
        
        result = re.sub(r'ide$', 'ayd', result)
        
        # Handle anything ending with -one
        
        result = re.sub(r'one$', 'own', result)
        
        # Handle anything ending with -e
        
        result = re.sub(r'e$', '', result)
        
        # Anything with thion should be tayon
        
        result = re.sub(r'thion', 'tayon', result)
        
        # Handle th specifically
        
        result = result.replace('th', 't')
        
        # Sub out for double letters in a row
        
        result = re.sub(r'([^\d])\1', r'\1', result)
    
        # Anything with ein turns into eene
        result = result.replace('ein', 'in')
        
        # Anything with 'o[vowel][!consonant]' has a y in middle
        
        result = re.sub(r'o([aeiou])(?![aeiou])', r'oy\1', result)
        
        # Anything with 'io' has a y in middle
        
        result = re.sub(r'i([aeiou])', r'iy\1', result)
        
        # Anything with 'eu' should be u in middle
        
        result = result.replace('eu', 'u')
        
        # 'ee' -> 'i', 'oo' -> 'u'
        result = result.replace('ee', 'i')
        result = result.replace('oo', 'u')
        

        
        # Step 3: Apply Spanish-specific patterns
        # if source_lang in ['spanish', 'auto']:
        #     for pattern, replacement in self.spanish_patterns:
        #         result = re.sub(pattern, replacement, result)
        
        # Step 6: Apply vowel shifts
        # for eng, fil in self.vowel_shifts.items():
        #     result = result.replace(eng, fil)
        
        # Step 4: Apply English-specific patterns
        if source_lang in ['english', 'auto']:
            for pattern, replacement in self.english_patterns:
                result = re.sub(pattern, replacement, result)
        
        # Step 5: Handle double consonants (reduce to single in most cases)
        # But preserve important doubles like 'dd' in 'adda'
        for consonant in 'bdfghjklmnpqrstvwxz':
            if consonant + consonant + consonant in result:
                result = result.replace(consonant + consonant + consonant, consonant)
        
        
        return result
    
    def add_affixes(self, root: str, affix_type: str = 'mag') -> str:
        """
        Add Filipino verbal affixes (Section 3.2.2 - Verb focus and aspect)
        
        Affix types:
        - mag: actor focus (mag-aral = to study)
        - nag: completed actor focus
        - um: actor focus infix
        - in: object focus suffix
        - pag: nominalizer
        - ka: superlative/recent completion
        """
        affixes = {
            'mag': f'mag-{root}',           # magkain ka rinaldo ng itlog
            'nag': f'nag-{root}',           # nagkain ni rinaldo ang kanyang talong
            'um': self._insert_um(root),    # kumain ka lang clarence
            'in': f'{root}-in',             # kainin mo ang sushi roll clive
            'an': f'{root}-an',             
            'pag': f'pag-{root}',           # ang pagkain ni Roan ay maliit
            'i': f'i-{root}',               
            'ka': f'ka-{root}',             # kakain lang ni rinaldo ang buldak
            'pang': f'pang-{root}',
        }
        
        return affixes.get(affix_type, root)
    
    def _insert_um(self, word: str) -> str:
        """Insert 'um' infix after first consonant or at beginning if starts with vowel"""
        vowels = 'aeiou'
        
        # If starts with vowel, just prefix
        if word[0] in vowels:
            return f'um{word}'
        
        # Find first vowel and insert 'um' before it
        for i, char in enumerate(word):
            if char in vowels:
                return word[:i] + 'um' + word[i:]
        
        return f'um{word}'
    
    def apply_reduplication(self, word: str, pattern: str = 'full', prefix: str = '') -> str:
        """
        Apply Filipino reduplication patterns (Section 2.3.5)
        
        For recent completion (ka-): reduplicate first syllable of ROOT, not prefix
        Example: ka-gising -> ka-gi-gising (not ka-ka-gising)
        """
        if pattern == 'full':
            return f'{word}-{word}'
        
        elif pattern == 'partial':
            # Reduplicate first CV (consonant-vowel)
            vowels = 'aeiou'
            
            # Skip prefix if provided
            root = word
            if prefix and word.startswith(prefix):
                root = word[len(prefix):]
            
            # Find first CV pattern in root
            for i in range(len(root) - 1):
                if root[i] not in vowels and root[i+1] in vowels:
                    reduplication = root[:i+1]
                    if prefix:
                        return f'{prefix}{reduplication}-{root}'
                    else:
                        return f'{reduplication}-{root}'
            
            # Fallback
            return f'{word[:2]}-{word}'
        
        elif pattern == 'recent':
            # ka- + CV-reduplication (Section 2.3.5)
            # kagigising = ka-gi-gising (reduplicate first syllable of "gising")
            vowels = 'aeiou'
            
            for i in range(len(word) - 1):
                if word[i] not in vowels and word[i+1] in vowels:
                    first_syllable = word[:i+2]  # Get CV
                    return f'ka{first_syllable}-{word}'
            
            return f'ka{word[:2]}-{word}'
        
        return word
    
    def apply_ligature(self, adjective: str, noun: str) -> str:
        """
        Apply Filipino ligature rules (Section 2.5.1)
        
        Rules:
        - If adjective ends in vowel: add -ng
        - If adjective ends in 'n': add -g  
        - If adjective ends in other consonant: add ' na'
        """
        last_char = adjective[-1].lower()
        
        if last_char in 'aeiou':
            # Vowel ending: add -ng
            return f'{adjective}ng {noun}'
        elif last_char == 'n':
            # Ends in 'n': add -g
            return f'{adjective}g {noun}'
        else:
            # Other consonant: use 'na'
            return f'{adjective} na {noun}'
    
    def check_particle_agreement(self, prev_word: str, particle: str) -> str:
        """
        Check din/rin agreement (Section 2.3.3)
        
        Rule: Use 'rin' after vowels/w/y, 'din' after consonants
        """
        if not prev_word:
            return particle
        
        last_char = prev_word.rstrip('.,!?;:')[-1].lower()
        
        if last_char in 'aeiouy' or last_char == 'w':
            return 'rin'
        else:
            return 'din'
    
    def parse_cfg_rule(self, rule: str) -> Tuple[str, List[str]]:
        """Parse CFG rule: S -> NP VP"""
        parts = rule.split('->')
        if len(parts) != 2:
            raise ValueError(f"Invalid CFG rule: {rule}")
        
        left = parts[0].strip()
        right = parts[1].strip().split()
        
        return left, right
    
    def demonstrate_naturalization(self):
        """Demonstrate Filipino naturalization with real examples from the paper"""
        print("=" * 70)
        print("FILIPINO NATURALIZATION PARSER")
        print("Based on: Ang, Chua, Tabe (2025) - Analysis of Filipino Spelling")
        print("=" * 70)
        
        # Spanish loanwords (Section 2.1.4)
        print("\n1. SPANISH LOANWORDS -> FILIPINO:")
        spanish_words = [
            ('teléfono', 'telepono'),
            ('verde', 'berde'),
            ('chocolate', 'tsokolate'),
            ('jeep', 'dyip'),
            ('zipon', 'sipon'),
            ('baño', 'banyo'),
            ('examen', 'eksamen'),
            ('centro', 'sentro'),
            ('como está', 'kumusta'),
            ('gobierno', 'gobyerno'),
            ('educacion', 'edukasyon'),
        ]
        for spanish, expected in spanish_words:
            result = self.apply_phonological_rules(spanish, 'spanish')
            match = "✓" if result == expected else "✗"
            print(f"   {match} {spanish:15} → {result:15} (expected: {expected})")
        
        # English loanwords
        print("\n2. ENGLISH LOANWORDS -> FILIPINO:")
        english_words = [
            ('computer', 'kompyuter'),
            ('facebook', 'peysbok'),
            ('television', 'telebisyon'),
            ('education', 'edukasyon'),
        ]
        for english, expected in english_words:
            result = self.apply_phonological_rules(english, 'english')
            print(f"   {english:15} → {result}")
        
        # Affixation (Section 3.2.2)
        print("\n3. AFFIXATION (Verb Focus):")
        examples = [
            ('aral', 'mag', 'mag-aral (to study)'),
            ('aral', 'nag', 'nag-aral (studied)'),
            ('kain', 'um', 'kumain (ate)'),
            ('sulat', 'in', 'sulat-in (to be written)'),
        ]
        for root, affix, description in examples:
            result = self.add_affixes(root, affix)
            print(f"   {root} + {affix:3} → {result:15} ({description})")
        
        # Reduplication for recent completion (Section 2.3.5)
        print("\n4. REDUPLICATION (Recent Completion - ka-):")
        print("   Rule: ka- + reduplicate first syllable of ROOT")
        examples = [
            ('gising', 'kagigising', 'just woke up'),
            ('sulat', 'kasusulat', 'just wrote'),
            ('tapos', 'katatapos', 'just finished'),
        ]
        for root, expected, meaning in examples:
            result = self.apply_reduplication(root, 'recent')
            match = "✓" if result.replace('-', '') == expected else "✗"
            print(f"   {match} {root:10} → {result:15} (expected: {expected}) = {meaning}")
        
        # Ligatures (Section 2.5.1)
        print("\n5. LIGATURES (Adjective-Noun):")
        examples = [
            ('maganda', 'bahay', 'magandang bahay', 'beautiful house'),
            ('malinis', 'silid', 'malinis na silid', 'clean room'),
            ('magaan', 'aklat', 'magaang aklat', 'light book'),
        ]
        for adj, noun, expected, meaning in examples:
            result = self.apply_ligature(adj, noun)
            match = "✓" if result == expected else "✗"
            print(f"   {match} {adj} + {noun:6} → {result:20} ({meaning})")
        
        # Particle agreement (Section 2.3.3)
        print("\n6. ENCLITIC PARTICLES (din/rin):")
        print("   Rule: 'rin' after vowels/w/y, 'din' after consonants")
        examples = [
            ('ako', 'rin', 'ako rin'),
            ('Juan', 'din', 'Juan din'),
            ('mahirap', 'din', 'mahirap din'),
            ('wala pa', 'rin', 'wala pa rin'),
        ]
        for prev, expected_particle, phrase in examples:
            result = self.check_particle_agreement(prev, 'din/rin')
            print(f"   {prev:10} + {result:4} → {phrase}")
        
        print("\n" + "=" * 70)
        
    def evaluate_parser(self, drug_table):
        TP = FP = 0
        for drug, expected_list in drug_table.items():
            result = self.apply_phonological_rules(drug, 'english')
            if result in expected_list:
                TP += 1
            else:
                FP += 1
        accuracy = TP / (TP + FP) if (TP + FP) else 0
        precision = TP / (TP + FP) if (TP + FP) else 0
        return {'TP': TP, 'FP': FP, 'Accuracy': accuracy, 'Precision': precision}



# Example usage
if __name__ == "__main__":
    parser = FilipinoCFGParser()
    
    # Run comprehensive demonstrations
    parser.demonstrate_naturalization()
    
    # Interactive examples
    print("\n" + "=" * 70)
    print("INTERACTIVE EXAMPLES")
    print("=" * 70)
    
    # print("\nSpanish words:")
    # for word in ['gobierno', 'educacion', 'teléfono', 'baño']:
    #     print(f"  {word:15} → {parser.apply_phonological_rules(word, 'spanish')}")
    
    # print("\nEnglish words:")
    # for word in ['computer', 'facebook', 'television']:
    #     print(f"  {word:15} → {parser.apply_phonological_rules(word, 'english')}")
        
    # Dictionary of drug names -> list of acceptable outputs
    drug_table = {
        'Acetaminophen': ['asetaminofen'],
        'Acetylcysteine': ['asetilsisteen'],
        'Acyclovir': ['asayklobir'],
        'Albendazole': ['albendasol'],
        'Albuterol': ['albiyuterol'],
        'Alendronate': ['alendroneyt'],
        'Alprazolam': ['alprasolam'],
        'Amikacin': ['amikasin'],
        'Amlodipine': ['amlodipin'],
        'Amoxicillin': ['amoksisilin'],
        'Apixaban': ['apiksaban'],
        'Aripiprazole': ['aripiprasol'],
        'Aspirin': ['aspirin'],
        'Atorvastatin': ['atorbastatin'],
        'Azithromycin': ['asitromaysin'],
        'Aztreonam': ['astreonam'],
        'Bacitracin': ['basitrasin'],
        'Baloxavir': ['baloksabir'],
        'Beclomethasone': ['beklometasown'],
        'Betamethasone': ['betametasown'],
        'Bisacodyl': ['bisakodil'],
        'Brompheniramine': ['brompeniramin', 'bromfeniramin'],
        'Budesonide': ['budesonayd'],
        'Bupropion': ['bupropyon', 'bupropiyon'],
        'Buspirone': ['buspirown'],
        'Calcium Carbonate': ['kalsiyum karboneyt'],
        'Captopril': ['kaptopril'],
        'Carbamazepine': ['karbamasepin'],
        'Carbimazole': ['karbimasol'],
        'Carvedilol': ['karbedilol'],
        'Celecoxib': ['selekoksib'],
        'Cephalexin': ['sefaleksin'],
        'Cetirizine': ['setirisin', 'sitirisin'],
        'Chloramphenicol': ['kloramfenikol'],
        'Chloroquine': ['klorokwin'],
        'Chlorpromazine': ['klorpromasin'],
        'Ciprofloxacin': ['siprofloksasin'],
        'Clindamycin': ['klindamaysin'],
        'Clobetasol': ['klobetasol'],
        'Clonazepam': ['klonasepam'],
        'Clopidogrel': ['klopidogrel'],
        'Clotrimazole': ['klotrimasol', 'klotrimasowl'],
        'Clozapine': ['klosapin'],
        'Codeine': ['kodeyn', 'kowdeyn', 'kodin', 'kowdin', 'kowdeen', 'kodeen'],
        'Cyanocobalamin': ['sayanokobalamin'],
        'Cyclosporine': ['sayklosporin'],
        'Cyproheptadine': ['sayproheptadin'],
        'Daptomycin': ['daptomaysin'],
        'Desloratadine': ['desloratadin'],
        'Dexamethasone': ['deksametasown'],
        'Dextromethorphan': ['dekstrometorfan', 'dektrometorpan'],
        'Diazepam': ['diyasepam'],
        'Diclofenac': ['diklofenak'],
        'Dimenhydrinate': ['dimenhaydrineyt'],
        'Diphenhydramine': ['daypenhaydramin'],
        'Dolutegravir': ['dolutegrabir'],
        'Domperidone': ['domperidown'],
        'Doxycycline': ['doksisayklin'],
        'Duloxetine': ['duloksetin'],
        'Efavirenz': ['efabirens'],
        'Enalapril': ['enalapril'],
        'Ertapenem': ['ertapenem'],
        'Erythromycin': ['eritromaysin'],
        'Escitalopram': ['eskitalopram'],
        'Esomeprazole': ['esomeprasol', 'esomeprasowl'],
        'Ethambutol': ['etambutol'],
        'Famciclovir': ['famsiklobir'],
        'Fentanyl': ['fentanil'],
        'Ferrous Sulfate': ['peryus sulpeyt'],
        'Fexofenadine': ['feksofenadin'],
        'Finasteride': ['finasterayd'],
        'Fluconazole': ['flukonasol', 'flukonasowl'],
        'Fluoxetine': ['fluoksetin'],
        'Fluticasone': ['flutikasown'],
        'Folic Acid': ['folik asid'],
        'Formoterol': ['formoterol'],
        'Furosemide': ['furosemayd'],
        'Gabapentin': ['gabapentin'],
        'Gentamicin': ['hentamaysin'],
        'Glimepiride': ['glaymepirayd'],
        'Guaifenesin': ['gwaypenesin', 'gwafenesin'],
        'Haloperidol': ['haloperidol'],
        'Heparin': ['heparin'],
        'Hydrochlorothiazide': ['haydroklorotayasayd'],
        'Hydrocortisone': ['haydrokortisown'],
        'Hydroxychloroquine': ['haydroksiklorokwin'],
        'Hydroxyzine': ['haydroksisin'],
        'Ibuprofen': ['aaybuprofen', 'aybuprofen', 'aaybupropen', 'aybupropen'],
        'Imipenem': ['imipenem'],
        'Insulin': ['insulin'],
        'Interferon': ['interferon'],
        'Ipratropium': ['ipratropiyum'],
        'Isoniazid': ['aysoneyasid'],
        'Ivermectin': ['aybermek­tin'],
        'Ketoconazole': ['ketokonasol', 'ketonasowl'],
        'Ketorolac': ['ketorolak'],
        'Lamivudine': ['lamibudin'],
        'Lamotrigine': ['lamotridyin'],
        'Levamisole': ['lebamisol'],
        'Levetiracetam': ['lebetirasetam'],
        'Levocetirizine': ['lebosetirisin'],
        'Levofloxacin': ['lebofloksasin'],
        'Levothyroxine': ['lebotayroksin'],
        'Linezolid': ['linesolid'],
        'Lisinopril': ['lisinopril'],
        'Lithium': ['litiyum'],
        'Loperamide': ['loperamayd'],
        'Loratadine': ['loratadin'],
        'Lorazepam': ['lorasepam'],
        'Losartan': ['losartan'],
        'Magnesium Oxide': ['magnesiyum oksayd'],
        'Malathion': ['malatayon'],
        'Mebendazole': ['mebendasol'],
        'Meclizine': ['meklisin'],
        'Meloxicam': ['meloksikam'],
        'Meropenem': ['meropenem'],
        'Metformin': ['metformin'],
        'Methimazole': ['metimasol'],
        'Methylprednisolone': ['metilprednisolown'],
        'Metoclopramide': ['metoklopramayd'],
        'Metoprolol': ['metoprolol'],
        'Metronidazole': ['metronidasol', 'metronidasowl'],
        'Midazolam': ['midasolam'],
        'Minocycline': ['minosayklin'],
        'Mirtazapine': ['mirtasapin'],
        'Molnupiravir': ['molnupirabir'],
        'Mometasone': ['mometasown'],
        'Montelukast': ['montelukast'],
        'Morphine': ['morfin', 'morpin'],
        'Naproxen': ['naproksen'],
        'Neomycin': ['neomaysin'],
        'Nitrofurantoin': ['nitropurantoyin', 'nitrofurantoyin'],
        'Ofloxacin': ['ofloksasin'],
        'Olanzapine': ['olansapin'],
        'Omeprazole': ['omeprasol'],
        'Ondansetron': ['ondansetron'],
        'Oseltamivir': ['oseltamibir'],
        'Oxycodone': ['oksikodown'],
        'Oxymetazoline': ['oksimetasolin'],
        'Pantoprazole': ['pantoprasol', 'pantoprasowl'],
        'Paracetamol': ['parasetamol'],
        'Paroxetine': ['paroksetin'],
        'Paxlovid': ['pakslobid'],
        'Permethrin': ['permetrin'],
        'Phenylephrine': ['fenilefrin'],
        'Phenytoin': ['penitoyn', 'fenitoyn'],
        'Polymyxin B': ['polimiksin bi'],
        'Praziquantel': ['prasikwantel'],
        'Prednisone': ['prednisown'],
        'Pregabalin': ['pregabalin'],
        'Prochlorperazine': ['proklorperasin'],
        'Promethazine': ['prometasin'],
        'Propranolol': ['propranolol'],
        'Pseudoephedrine': ['sudoyepedrin', 'sudoyefedrin'],
        'Psyllium': ['sillyum', 'silium', 'siliyum'],
        'Pyrantel Pamoate': ['payrantel pamoeyt'],
        'Pyrazinamide': ['payrasinamayd'],
        'Quetiapine': ['kwetiyapin'],
        'Ranitidine': ['ranitidin'],
        'Remdesivir': ['remdesibir'],
        'Rifampicin': ['rifampisin'],
        'Risperidone': ['risperidown'],
        'Rivaroxaban': ['ribaroksaban'],
        'Rosuvastatin': ['rosubastatin'],
        'Salbutamol': ['salbutamol'],
        'Salmeterol': ['salmeterol'],
        'Senna': ['sena'],
        'Sertraline': ['sertralin'],
        'Sildenafil': ['sildenafil'],
        'Simethicone': ['simetikown'],
        'Simvastatin': ['simbastatin'],
        'Sitagliptin': ['sitagliptin'],
        'Spironolactone': ['spironolaktown'],
        'Sulfamethoxazole': ['sulfametoksasol'],
        'Tacrolimus': ['takrolimus'],
        'Tadalafil': ['tadalapil', 'tadalafil'],
        'Tamsulosin': ['tamsulosin'],
        'Tenofovir': ['tenofobir'],
        'Tetracycline': ['tetrasayklin'],
        'Theophylline': ['teofilin'],
        'Tigecycline': ['tigesayklin'],
        'Tiotropium': ['tiyotropiyum'],
        'Tobramycin': ['tobramaysin'],
        'Topiramate': ['topirameyt'],
        'Tramadol': ['tramadol'],
        'Trazodone': ['trasodown'],
        'Triamcinolone': ['triyamsinolown'],
        'Trimethoprim': ['trimetoprim'],
        'Valacyclovir': ['balasayklobir'],
        'Valproic Acid': ['balproyik asid'],
        'Valsartan': ['balsartan'],
        'Vancomycin': ['bankomaysin'],
        'Venlafaxine': ['benlafaksin'],
        'Vitamin C': ['bitamin si'],
        'Vitamin D': ['bitamin di'],
        'Warfarin': ['warfarin'],
        'Xylometazoline': ['silometasolin'],
        'Zanamivir': ['sanamibir'],
        'Zinc Sulfate': ['sink sulpeyt', 'sink sulfeyt'],
    }
    

    # Testing loop
    for drug, expected_list in drug_table.items():
        result = parser.apply_phonological_rules(drug, 'english')  # your parser function
        match = '✓' if result in expected_list else '✗'
        print(f"{match} {drug:25} → {result:25} (expected: {', '.join(expected_list)})")

    metrics = parser.evaluate_parser(drug_table)
    print(metrics)
    
    # print("\nWith affixes:")
    # print(f"  luto + mag → {parser.add_affixes('luto', 'mag')}")
    # print(f"  kain + um  → {parser.add_affixes('kain', 'um')}")
    
    # print("\nLigatures:")
    # print(f"  maganda + bahay → {parser.apply_ligature('maganda', 'bahay')}")
    # print(f"  malinis + silid → {parser.apply_ligature('malinis', 'silid')}")
    
    print("\n" + "=" * 70)